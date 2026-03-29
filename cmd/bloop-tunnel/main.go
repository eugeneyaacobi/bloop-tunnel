package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"bloop-tunnel/internal/auth"
	"bloop-tunnel/internal/client"
	clientsetup "bloop-tunnel/internal/client/setup"
	"bloop-tunnel/internal/config"
	"bloop-tunnel/internal/logging"
	"bloop-tunnel/pkg/version"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "setup" {
		os.Exit(runSetup(args[1:]))
	}

	configPath := flag.String("config", "", "Path to client config (optional when using environment variables)")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("bloop-tunnel %s (%s) %s\n", version.Version, version.Commit, version.Date)
		return
	}

	cfg, err := config.LoadClientConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load client config: %v\n", err)
		os.Exit(1)
	}

	token := auth.ResolveClientToken(cfg.AuthToken, cfg.AuthTokenEnv)
	logger := logging.New(cfg.Logging.Level)
	logger.Info("client starting", "relay_url", cfg.RelayURL, "tunnel_count", len(cfg.Tunnels), "has_token", token != "")

	if cfg.ControlPlaneURL != "" {
		enrollmentToken := cfg.EnrollmentToken
		if enrollmentToken == "" && cfg.EnrollmentTokenEnv != "" {
			enrollmentToken = os.Getenv(cfg.EnrollmentTokenEnv)
		}
		if enrollmentToken != "" {
			installationID, ingestToken, err := client.EnrollRuntime(context.Background(), cfg.ControlPlaneURL, enrollmentToken, "default-client")
			if err != nil {
				fmt.Fprintf(os.Stderr, "enroll runtime: %v\n", err)
				os.Exit(1)
			}
			logger.Info("runtime enrollment succeeded", "installation_id", installationID)
			os.Setenv("BLOOP_RUNTIME_INSTALLATION_ID", installationID)
			os.Setenv("BLOOP_RUNTIME_INGEST_BEARER", ingestToken)
		}
	}

	session, err := client.Connect(context.Background(), cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to relay: %v\n", err)
		os.Exit(1)
	}
	defer session.Transport.Close()

	if err := validateProtectedTunnels(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "invalid tunnel config: %v\n", err)
		os.Exit(1)
	}

	if err := session.RegisterTunnels(); err != nil {
		fmt.Fprintf(os.Stderr, "register tunnels: %v\n", err)
		os.Exit(1)
	}

	logger.Info("client registered tunnels successfully")
	registered := make([]string, 0, len(cfg.Tunnels))
	for _, t := range cfg.Tunnels {
		registered = append(registered, t.Name)
	}
	go func() {
		if err := client.StartRuntimeIngestLoop(context.Background(), cfg, registered); err != nil {
			logger.Error("runtime ingest loop failed", "error", err.Error())
		}
	}()
	if err := session.RunWithReconnect(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "client session ended: %v\n", err)
		os.Exit(1)
	}
}

func runSetup(args []string) int {
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := fs.String("config", "deploy/examples/client.example.yaml", "Write YAML scaffold to this path, or '-' for stdout")
	output := fs.String("output", string(clientsetup.OutputYAML), "Output mode: yaml|env-file|compose-block")
	nonInteractive := fs.Bool("non-interactive", false, "Generate scaffold output without prompts")
	controlPlaneURL := fs.String("control-plane-url", config.DefaultControlPlaneURL, "Control plane URL to embed in generated output")
	relayURL := fs.String("relay-url", config.DefaultRelayURL, "Relay URL to embed in generated output")
	authTokenEnv := fs.String("auth-token-env", "BLOOP_CLIENT_TOKEN", "Environment variable name holding the relay auth token")
	enrollmentTokenEnv := fs.String("enrollment-token-env", "", "Environment variable name holding the enrollment token")
	discoverDocker := fs.Bool("discover-docker", false, "Offer opt-in Docker service discovery during interactive setup")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	err := clientsetup.Run(os.Stdout, os.Stderr, clientsetup.Options{
		ConfigPath:         *configPath,
		OutputMode:         clientsetup.OutputMode(*output),
		NonInteractive:     *nonInteractive,
		ControlPlaneURL:    *controlPlaneURL,
		RelayURL:           *relayURL,
		AuthTokenEnv:       *authTokenEnv,
		EnrollmentTokenEnv: *enrollmentTokenEnv,
		DiscoverDocker:     *discoverDocker,
		Stdin:              os.Stdin,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup: %v\n", err)
		return 1
	}
	return 0
}

func validateProtectedTunnels(cfg *config.ClientConfig) error {
	for _, t := range cfg.Tunnels {
		switch t.Access {
		case "basic_auth":
			if t.BasicAuth.Username == "" {
				return fmt.Errorf("tunnel %q missing basic auth username", t.Name)
			}
			if t.BasicAuth.Password == "" && t.BasicAuth.PasswordEnv == "" {
				return fmt.Errorf("tunnel %q missing basic auth password or password_env", t.Name)
			}
		case "token_protected":
			if t.Token == "" && t.TokenEnv == "" {
				return fmt.Errorf("tunnel %q missing token or token_env", t.Name)
			}
		case "public", "":
		default:
			return errors.New("unsupported access mode for tunnel " + t.Name)
		}
	}
	return nil
}
