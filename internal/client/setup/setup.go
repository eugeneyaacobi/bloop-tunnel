package setup

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"bloop-tunnel/internal/client/dockerdiscover"
	"bloop-tunnel/internal/config"
	"gopkg.in/yaml.v3"
)

type OutputMode string

const (
	OutputYAML         OutputMode = "yaml"
	OutputEnvFile      OutputMode = "env-file"
	OutputComposeBlock OutputMode = "compose-block"
)

type Options struct {
	ConfigPath         string
	OutputMode         OutputMode
	NonInteractive     bool
	ControlPlaneURL    string
	RelayURL           string
	AuthTokenEnv       string
	EnrollmentTokenEnv string
	DiscoverDocker     bool
	Stdin              io.Reader
	Discoverer         dockerdiscover.Discoverer
}

func Run(stdout, stderr io.Writer, opts Options) error {
	if opts.OutputMode == "" {
		opts.OutputMode = OutputYAML
	}
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}

	cfg, err := initialConfig(opts)
	if err != nil {
		return err
	}
	if opts.NonInteractive {
		return writeOutput(stdout, opts, cfg)
	}
	if err := runInteractive(context.Background(), stdout, stderr, opts.Stdin, stderr, &cfg, opts); err != nil {
		return err
	}
	return writeOutput(stdout, opts, cfg)
}

func initialConfig(opts Options) (config.ClientConfig, error) {
	cfg := config.DefaultClientConfig()
	if opts.ConfigPath != "" && opts.ConfigPath != "-" {
		loaded, err := loadExistingConfig(opts.ConfigPath)
		if err != nil {
			return config.ClientConfig{}, err
		}
		cfg = loaded
	}
	if opts.ControlPlaneURL != "" {
		cfg.ControlPlaneURL = opts.ControlPlaneURL
	}
	if opts.RelayURL != "" {
		cfg.RelayURL = opts.RelayURL
	}
	if opts.AuthTokenEnv != "" {
		cfg.AuthTokenEnv = opts.AuthTokenEnv
	} else if cfg.AuthTokenEnv == "" {
		cfg.AuthTokenEnv = "BLOOP_CLIENT_TOKEN"
	}
	if opts.EnrollmentTokenEnv != "" {
		cfg.EnrollmentTokenEnv = opts.EnrollmentTokenEnv
	}
	if len(cfg.Tunnels) == 0 {
		cfg.Tunnels = []config.TunnelConfig{{
			Name:      "app",
			Hostname:  "app.bloop.to",
			LocalAddr: "127.0.0.1:3000",
			Access:    "public",
		}}
	}
	return cfg, nil
}

func runInteractive(ctx context.Context, _ io.Writer, stderr io.Writer, in io.Reader, promptOut io.Writer, cfg *config.ClientConfig, opts Options) error {
	reader := bufio.NewReader(in)
	fmt.Fprintln(promptOut, "Interactive bloop-tunnel setup")
	fmt.Fprintln(promptOut, "Press enter to accept defaults. Docker discovery is always opt-in.")

	cfg.ControlPlaneURL = prompt(reader, promptOut, "Control plane URL", cfg.ControlPlaneURL)
	cfg.RelayURL = prompt(reader, promptOut, "Relay URL", cfg.RelayURL)
	cfg.AuthTokenEnv = prompt(reader, promptOut, "Relay auth token environment variable", defaultString(cfg.AuthTokenEnv, "BLOOP_CLIENT_TOKEN"))
	cfg.EnrollmentTokenEnv = prompt(reader, promptOut, "Enrollment token environment variable (optional)", cfg.EnrollmentTokenEnv)
	cfg.Reconnect.InitialDelayMs = promptInt(reader, promptOut, "Reconnect initial delay ms", cfg.Reconnect.InitialDelayMs)
	cfg.Reconnect.MaxDelayMs = promptInt(reader, promptOut, "Reconnect max delay ms", cfg.Reconnect.MaxDelayMs)

	if opts.DiscoverDocker && confirm(reader, promptOut, "Inspect local Docker services before editing tunnels?", true) {
		if err := applyDockerDiscovery(ctx, reader, promptOut, cfg, opts); err != nil {
			fmt.Fprintf(stderr, "docker discovery unavailable: %v\n", err)
		}
	}

	for {
		fmt.Fprintln(promptOut, "Current tunnels:")
		for i, tunnel := range cfg.Tunnels {
			fmt.Fprintf(promptOut, "  [%d] %s -> %s (%s)\n", i+1, tunnel.Name, tunnel.LocalAddr, defaultString(tunnel.Access, "public"))
		}
		action := strings.ToLower(prompt(reader, promptOut, "Tunnel action [edit/add/remove/done]", "done"))
		switch action {
		case "", "done", "d":
			if len(cfg.Tunnels) == 0 {
				fmt.Fprintln(promptOut, "You need at least one tunnel.")
				continue
			}
			if err := validateConfig(*cfg); err != nil {
				fmt.Fprintf(promptOut, "Tunnel validation failed: %v\n", err)
				continue
			}
			return nil
		case "edit", "e":
			idx := promptIndex(reader, promptOut, "Tunnel number to edit", len(cfg.Tunnels))
			cfg.Tunnels[idx] = editTunnel(reader, promptOut, cfg.Tunnels[idx])
		case "add", "a":
			cfg.Tunnels = append(cfg.Tunnels, editTunnel(reader, promptOut, config.TunnelConfig{Access: "public"}))
		case "remove", "r":
			if len(cfg.Tunnels) == 0 {
				fmt.Fprintln(promptOut, "No tunnels to remove.")
				continue
			}
			idx := promptIndex(reader, promptOut, "Tunnel number to remove", len(cfg.Tunnels))
			cfg.Tunnels = append(cfg.Tunnels[:idx], cfg.Tunnels[idx+1:]...)
		default:
			fmt.Fprintln(promptOut, "Unknown action. Use edit, add, remove, or done.")
		}
	}
}

func loadExistingConfig(path string) (config.ClientConfig, error) {
	cfg := config.DefaultClientConfig()
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return config.ClientConfig{}, err
	}
	var fileCfg config.ClientConfig
	if err := yaml.Unmarshal(payload, &fileCfg); err != nil {
		return config.ClientConfig{}, err
	}
	if fileCfg.ControlPlaneURL != "" {
		cfg.ControlPlaneURL = fileCfg.ControlPlaneURL
	}
	if fileCfg.RelayURL != "" {
		cfg.RelayURL = fileCfg.RelayURL
	}
	if fileCfg.AuthToken != "" {
		cfg.AuthToken = fileCfg.AuthToken
	}
	if fileCfg.AuthTokenEnv != "" {
		cfg.AuthTokenEnv = fileCfg.AuthTokenEnv
	}
	if fileCfg.EnrollmentToken != "" {
		cfg.EnrollmentToken = fileCfg.EnrollmentToken
	}
	if fileCfg.EnrollmentTokenEnv != "" {
		cfg.EnrollmentTokenEnv = fileCfg.EnrollmentTokenEnv
	}
	if fileCfg.Reconnect.InitialDelayMs != 0 {
		cfg.Reconnect.InitialDelayMs = fileCfg.Reconnect.InitialDelayMs
	}
	if fileCfg.Reconnect.MaxDelayMs != 0 {
		cfg.Reconnect.MaxDelayMs = fileCfg.Reconnect.MaxDelayMs
	}
	if fileCfg.Logging.Level != "" {
		cfg.Logging.Level = fileCfg.Logging.Level
	}
	if fileCfg.Logging.Format != "" {
		cfg.Logging.Format = fileCfg.Logging.Format
	}
	if len(fileCfg.Tunnels) > 0 {
		cfg.Tunnels = fileCfg.Tunnels
	}
	return cfg, nil
}

func writeOutput(stdout io.Writer, opts Options, cfg config.ClientConfig) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}
	switch opts.OutputMode {
	case OutputYAML:
		payload, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		if opts.ConfigPath == "" || opts.ConfigPath == "-" {
			_, err = stdout.Write(payload)
			return err
		}
		if err := os.WriteFile(opts.ConfigPath, payload, 0o600); err != nil {
			return err
		}
		_, err = fmt.Fprintf(stdout, "wrote setup output to %s\n", opts.ConfigPath)
		return err
	case OutputEnvFile:
		_, err := io.WriteString(stdout, renderEnv(cfg, false))
		return err
	case OutputComposeBlock:
		_, err := io.WriteString(stdout, renderEnv(cfg, true))
		return err
	default:
		return fmt.Errorf("unsupported output mode %q", opts.OutputMode)
	}
}

func validateConfig(cfg config.ClientConfig) error {
	if strings.TrimSpace(cfg.ControlPlaneURL) == "" {
		return fmt.Errorf("control_plane_url is required")
	}
	if strings.TrimSpace(cfg.RelayURL) == "" {
		return fmt.Errorf("relay_url is required")
	}
	if len(cfg.Tunnels) == 0 {
		return fmt.Errorf("at least one tunnel is required")
	}
	for i, tunnel := range cfg.Tunnels {
		if strings.TrimSpace(tunnel.Name) == "" || strings.TrimSpace(tunnel.LocalAddr) == "" {
			return fmt.Errorf("tunnels[%d] requires name and local_addr", i)
		}
	}
	return nil
}

func renderEnv(cfg config.ClientConfig, compose bool) string {
	lines := []string{
		fmt.Sprintf("BLOOP_CONTROL_PLANE_URL=%s", cfg.ControlPlaneURL),
		fmt.Sprintf("BLOOP_RELAY_URL=%s", cfg.RelayURL),
	}
	if cfg.AuthTokenEnv != "" {
		lines = append(lines, fmt.Sprintf("BLOOP_AUTH_TOKEN_ENV=%s", cfg.AuthTokenEnv))
	}
	if cfg.EnrollmentTokenEnv != "" {
		lines = append(lines, fmt.Sprintf("BLOOP_ENROLLMENT_TOKEN_ENV=%s", cfg.EnrollmentTokenEnv))
	}
	for i, tunnel := range cfg.Tunnels {
		prefix := fmt.Sprintf("BLOOP_TUNNELS_%d_", i)
		lines = append(lines,
			prefix+"NAME="+tunnel.Name,
			prefix+"HOSTNAME="+tunnel.Hostname,
			prefix+"LOCAL_ADDR="+tunnel.LocalAddr,
			prefix+"ACCESS="+defaultString(tunnel.Access, "public"),
		)
		if tunnel.BasicAuth.Username != "" {
			lines = append(lines, prefix+"BASIC_AUTH_USERNAME="+tunnel.BasicAuth.Username)
		}
		if tunnel.BasicAuth.PasswordEnv != "" {
			lines = append(lines, prefix+"BASIC_AUTH_PASSWORD_ENV="+tunnel.BasicAuth.PasswordEnv)
		}
		if tunnel.TokenEnv != "" {
			lines = append(lines, prefix+"TOKEN_ENV="+tunnel.TokenEnv)
		}
	}
	if !compose {
		return strings.Join(lines, "\n") + "\n"
	}
	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		lines[i] = fmt.Sprintf("      %s: %s", parts[0], parts[1])
	}
	return "environment:\n" + strings.Join(lines, "\n") + "\n"
}

func prompt(reader *bufio.Reader, out io.Writer, label, current string) string {
	if current != "" {
		fmt.Fprintf(out, "%s [%s]: ", label, current)
	} else {
		fmt.Fprintf(out, "%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return current
	}
	return line
}

func promptInt(reader *bufio.Reader, out io.Writer, label string, current int) int {
	for {
		value := prompt(reader, out, label, strconv.Itoa(current))
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
		fmt.Fprintf(out, "Enter a whole number.\n")
	}
}

func promptIndex(reader *bufio.Reader, out io.Writer, label string, max int) int {
	for {
		value := prompt(reader, out, label, "1")
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 1 && parsed <= max {
			return parsed - 1
		}
		fmt.Fprintf(out, "Choose a number between 1 and %d.\n", max)
	}
}

func editTunnel(reader *bufio.Reader, out io.Writer, tunnel config.TunnelConfig) config.TunnelConfig {
	tunnel.Name = prompt(reader, out, "Tunnel name", tunnel.Name)
	tunnel.Hostname = prompt(reader, out, "Hostname (optional)", tunnel.Hostname)
	tunnel.LocalAddr = prompt(reader, out, "Local address", tunnel.LocalAddr)
	tunnel.Access = strings.ToLower(prompt(reader, out, "Access [public/basic_auth/token_protected]", defaultString(tunnel.Access, "public")))
	switch tunnel.Access {
	case "basic_auth":
		tunnel.BasicAuth.Username = prompt(reader, out, "Basic auth username", tunnel.BasicAuth.Username)
		tunnel.BasicAuth.PasswordEnv = prompt(reader, out, "Basic auth password env", tunnel.BasicAuth.PasswordEnv)
		tunnel.TokenEnv = ""
	case "token_protected":
		tunnel.TokenEnv = prompt(reader, out, "Token env", tunnel.TokenEnv)
		tunnel.BasicAuth = config.BasicAuthConfig{}
	default:
		tunnel.Access = "public"
		tunnel.BasicAuth = config.BasicAuthConfig{}
		tunnel.TokenEnv = ""
	}
	return tunnel
}

func confirm(reader *bufio.Reader, out io.Writer, label string, current bool) bool {
	def := "y/N"
	if current {
		def = "Y/n"
	}
	for {
		fmt.Fprintf(out, "%s [%s]: ", label, def)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			return current
		}
		if line == "y" || line == "yes" {
			return true
		}
		if line == "n" || line == "no" {
			return false
		}
	}
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
