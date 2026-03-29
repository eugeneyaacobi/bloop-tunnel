package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultClientConfigUsesProductionDefaults(t *testing.T) {
	cfg := DefaultClientConfig()
	if cfg.ControlPlaneURL != DefaultControlPlaneURL {
		t.Fatalf("control plane default = %q, want %q", cfg.ControlPlaneURL, DefaultControlPlaneURL)
	}
	if cfg.RelayURL != DefaultRelayURL {
		t.Fatalf("relay default = %q, want %q", cfg.RelayURL, DefaultRelayURL)
	}
}

func TestLoadClientConfigMergesFileAndEnvWithEnvPrecedence(t *testing.T) {
	t.Setenv(EnvControlPlaneURL, "https://env-api.bloop.to")
	t.Setenv(EnvRelayURL, "wss://env-relay.bloop.to/connect")
	t.Setenv(EnvAuthToken, "env-token")
	t.Setenv(EnvLoggingLevel, "debug")
	t.Setenv("BLOOP_TUNNELS_0_NAME", "env-app")
	t.Setenv("BLOOP_TUNNELS_0_LOCAL_ADDR", "host.docker.internal:3000")
	t.Setenv("BLOOP_TUNNELS_0_ACCESS", "public")

	dir := t.TempDir()
	path := filepath.Join(dir, "client.yaml")
	content := strings.TrimSpace(`
relay_url: wss://file-relay.bloop.to/connect
auth_token_env: FILE_TOKEN
control_plane_url: https://file-api.bloop.to
logging:
  level: info
tunnels:
  - name: file-app
    local_addr: 127.0.0.1:8080
    access: public
`) + "\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadClientConfig(path)
	if err != nil {
		t.Fatalf("LoadClientConfig returned error: %v", err)
	}
	if cfg.ControlPlaneURL != "https://env-api.bloop.to" {
		t.Fatalf("control plane = %q", cfg.ControlPlaneURL)
	}
	if cfg.RelayURL != "wss://env-relay.bloop.to/connect" {
		t.Fatalf("relay = %q", cfg.RelayURL)
	}
	if cfg.AuthToken != "env-token" {
		t.Fatalf("auth token = %q", cfg.AuthToken)
	}
	if cfg.Logging.Level != "debug" {
		t.Fatalf("logging level = %q", cfg.Logging.Level)
	}
	if len(cfg.Tunnels) != 1 || cfg.Tunnels[0].Name != "env-app" {
		t.Fatalf("tunnels = %#v", cfg.Tunnels)
	}
}

func TestLoadClientConfigEnvOnly(t *testing.T) {
	t.Setenv(EnvAuthToken, "env-token")
	t.Setenv("BLOOP_TUNNELS_0_NAME", "app")
	t.Setenv("BLOOP_TUNNELS_0_LOCAL_ADDR", "host.docker.internal:3000")
	t.Setenv("BLOOP_TUNNELS_1_NAME", "admin")
	t.Setenv("BLOOP_TUNNELS_1_LOCAL_ADDR", "host.docker.internal:4000")
	t.Setenv("BLOOP_TUNNELS_1_ACCESS", "basic_auth")
	t.Setenv("BLOOP_TUNNELS_1_BASIC_AUTH_USERNAME", "gene")
	t.Setenv("BLOOP_TUNNELS_1_BASIC_AUTH_PASSWORD_ENV", "BLOOP_ADMIN_PASSWORD")

	cfg, err := LoadClientConfig("")
	if err != nil {
		t.Fatalf("LoadClientConfig returned error: %v", err)
	}
	if cfg.ControlPlaneURL != DefaultControlPlaneURL {
		t.Fatalf("control plane = %q", cfg.ControlPlaneURL)
	}
	if len(cfg.Tunnels) != 2 {
		t.Fatalf("tunnel count = %d", len(cfg.Tunnels))
	}
	if cfg.Tunnels[1].BasicAuth.PasswordEnv != "BLOOP_ADMIN_PASSWORD" {
		t.Fatalf("password env = %q", cfg.Tunnels[1].BasicAuth.PasswordEnv)
	}
}

func TestDecodeTunnelEnvRequiresNameAndLocalAddr(t *testing.T) {
	_, err := DecodeTunnelEnv([]string{"BLOOP_TUNNELS_0_NAME=app"})
	if err == nil {
		t.Fatal("expected error")
	}
}
