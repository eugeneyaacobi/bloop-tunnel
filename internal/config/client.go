package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultControlPlaneURL     = "https://api.bloop.to"
	DefaultRelayURL            = "wss://relay.bloop.to/connect"
	DefaultReconnectInitialMS  = 1000
	DefaultReconnectMaxMS      = 30000
	DefaultLoggingLevel        = "info"
	DefaultLoggingFormat       = "json"
	EnvControlPlaneURL         = "BLOOP_CONTROL_PLANE_URL"
	EnvRelayURL                = "BLOOP_RELAY_URL"
	EnvAuthToken               = "BLOOP_AUTH_TOKEN"
	EnvAuthTokenEnv            = "BLOOP_AUTH_TOKEN_ENV"
	EnvEnrollmentToken         = "BLOOP_ENROLLMENT_TOKEN"
	EnvEnrollmentTokenEnv      = "BLOOP_ENROLLMENT_TOKEN_ENV"
	EnvReconnectInitialDelayMs = "BLOOP_RECONNECT_INITIAL_DELAY_MS"
	EnvReconnectMaxDelayMs     = "BLOOP_RECONNECT_MAX_DELAY_MS"
	EnvLoggingLevel            = "BLOOP_LOG_LEVEL"
	EnvLoggingFormat           = "BLOOP_LOG_FORMAT"
	EnvTunnelPrefix            = "BLOOP_TUNNELS_"
)

type ReconnectConfig struct {
	InitialDelayMs int `mapstructure:"initial_delay_ms" yaml:"initial_delay_ms"`
	MaxDelayMs     int `mapstructure:"max_delay_ms" yaml:"max_delay_ms"`
}

type BasicAuthConfig struct {
	Username    string `mapstructure:"username" yaml:"username,omitempty"`
	Password    string `mapstructure:"password" yaml:"password,omitempty"`
	PasswordEnv string `mapstructure:"password_env" yaml:"password_env,omitempty"`
}

type TunnelConfig struct {
	Name      string          `mapstructure:"name" yaml:"name"`
	Hostname  string          `mapstructure:"hostname" yaml:"hostname,omitempty"`
	LocalAddr string          `mapstructure:"local_addr" yaml:"local_addr"`
	Access    string          `mapstructure:"access" yaml:"access,omitempty"`
	Token     string          `mapstructure:"token" yaml:"token,omitempty"`
	TokenEnv  string          `mapstructure:"token_env" yaml:"token_env,omitempty"`
	BasicAuth BasicAuthConfig `mapstructure:"basic_auth" yaml:"basic_auth,omitempty"`
}

type ClientConfig struct {
	RelayURL           string          `mapstructure:"relay_url" yaml:"relay_url"`
	AuthToken          string          `mapstructure:"auth_token" yaml:"auth_token,omitempty"`
	AuthTokenEnv       string          `mapstructure:"auth_token_env" yaml:"auth_token_env,omitempty"`
	ControlPlaneURL    string          `mapstructure:"control_plane_url" yaml:"control_plane_url"`
	EnrollmentToken    string          `mapstructure:"enrollment_token" yaml:"enrollment_token,omitempty"`
	EnrollmentTokenEnv string          `mapstructure:"enrollment_token_env" yaml:"enrollment_token_env,omitempty"`
	Reconnect          ReconnectConfig `mapstructure:"reconnect" yaml:"reconnect"`
	Tunnels            []TunnelConfig  `mapstructure:"tunnels" yaml:"tunnels,omitempty"`
	Logging            Logging         `mapstructure:"logging" yaml:"logging"`
}

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		RelayURL:        DefaultRelayURL,
		ControlPlaneURL: DefaultControlPlaneURL,
		Reconnect: ReconnectConfig{
			InitialDelayMs: DefaultReconnectInitialMS,
			MaxDelayMs:     DefaultReconnectMaxMS,
		},
		Logging: Logging{
			Level:  DefaultLoggingLevel,
			Format: DefaultLoggingFormat,
		},
	}
}

func LoadClientConfig(path string) (*ClientConfig, error) {
	resolved, err := ResolveClientConfig(path)
	if err != nil {
		return nil, err
	}
	return &resolved, nil
}

func ResolveClientConfig(path string) (ClientConfig, error) {
	cfg := DefaultClientConfig()

	if path != "" {
		if err := mergeClientFile(path, &cfg); err != nil {
			return ClientConfig{}, err
		}
	}

	mergeClientEnv(os.Environ(), &cfg)
	if err := validateClientConfig(cfg); err != nil {
		return ClientConfig{}, err
	}
	return cfg, nil
}

func mergeClientFile(path string, cfg *ClientConfig) error {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFound) || errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	var fileCfg ClientConfig
	if err := v.Unmarshal(&fileCfg); err != nil {
		return err
	}

	mergeClientConfig(cfg, fileCfg)
	return nil
}

func mergeClientConfig(dst *ClientConfig, src ClientConfig) {
	if src.RelayURL != "" {
		dst.RelayURL = src.RelayURL
	}
	if src.AuthToken != "" {
		dst.AuthToken = src.AuthToken
	}
	if src.AuthTokenEnv != "" {
		dst.AuthTokenEnv = src.AuthTokenEnv
	}
	if src.ControlPlaneURL != "" {
		dst.ControlPlaneURL = src.ControlPlaneURL
	}
	if src.EnrollmentToken != "" {
		dst.EnrollmentToken = src.EnrollmentToken
	}
	if src.EnrollmentTokenEnv != "" {
		dst.EnrollmentTokenEnv = src.EnrollmentTokenEnv
	}
	if src.Reconnect.InitialDelayMs != 0 {
		dst.Reconnect.InitialDelayMs = src.Reconnect.InitialDelayMs
	}
	if src.Reconnect.MaxDelayMs != 0 {
		dst.Reconnect.MaxDelayMs = src.Reconnect.MaxDelayMs
	}
	if src.Logging.Level != "" {
		dst.Logging.Level = src.Logging.Level
	}
	if src.Logging.Format != "" {
		dst.Logging.Format = src.Logging.Format
	}
	if len(src.Tunnels) > 0 {
		dst.Tunnels = append([]TunnelConfig(nil), src.Tunnels...)
	}
}

func mergeClientEnv(env []string, cfg *ClientConfig) {
	lookup := envMap(env)
	if value := lookup[EnvRelayURL]; value != "" {
		cfg.RelayURL = value
	}
	if value := lookup[EnvAuthToken]; value != "" {
		cfg.AuthToken = value
	}
	if value := lookup[EnvAuthTokenEnv]; value != "" {
		cfg.AuthTokenEnv = value
	}
	if value := lookup[EnvControlPlaneURL]; value != "" {
		cfg.ControlPlaneURL = value
	}
	if value := lookup[EnvEnrollmentToken]; value != "" {
		cfg.EnrollmentToken = value
	}
	if value := lookup[EnvEnrollmentTokenEnv]; value != "" {
		cfg.EnrollmentTokenEnv = value
	}
	if value := lookup[EnvReconnectInitialDelayMs]; value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			cfg.Reconnect.InitialDelayMs = parsed
		}
	}
	if value := lookup[EnvReconnectMaxDelayMs]; value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			cfg.Reconnect.MaxDelayMs = parsed
		}
	}
	if value := lookup[EnvLoggingLevel]; value != "" {
		cfg.Logging.Level = value
	}
	if value := lookup[EnvLoggingFormat]; value != "" {
		cfg.Logging.Format = value
	}

	tunnels, err := DecodeTunnelEnv(env)
	if err == nil && len(tunnels) > 0 {
		cfg.Tunnels = tunnels
	}
}

func validateClientConfig(cfg ClientConfig) error {
	for i, tunnel := range cfg.Tunnels {
		if tunnel.Name == "" {
			return fmt.Errorf("tunnels[%d].name is required", i)
		}
		if tunnel.LocalAddr == "" {
			return fmt.Errorf("tunnels[%d].local_addr is required", i)
		}
	}
	return nil
}

func envMap(env []string) map[string]string {
	lookup := make(map[string]string, len(env))
	for _, item := range env {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			continue
		}
		lookup[parts[0]] = parts[1]
	}
	return lookup
}

func DecodeTunnelEnv(env []string) ([]TunnelConfig, error) {
	lookup := envMap(env)
	indexed := map[int]*TunnelConfig{}

	for key, value := range lookup {
		if !strings.HasPrefix(key, EnvTunnelPrefix) {
			continue
		}
		remainder := strings.TrimPrefix(key, EnvTunnelPrefix)
		parts := strings.Split(remainder, "_")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid tunnel env key %q", key)
		}
		idx, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid tunnel index in %q", key)
		}
		field := strings.Join(parts[1:], "_")
		tunnel := indexed[idx]
		if tunnel == nil {
			tunnel = &TunnelConfig{}
			indexed[idx] = tunnel
		}
		switch field {
		case "NAME":
			tunnel.Name = value
		case "HOSTNAME":
			tunnel.Hostname = value
		case "LOCAL_ADDR":
			tunnel.LocalAddr = value
		case "ACCESS":
			tunnel.Access = value
		case "TOKEN":
			tunnel.Token = value
		case "TOKEN_ENV":
			tunnel.TokenEnv = value
		case "BASIC_AUTH_USERNAME":
			tunnel.BasicAuth.Username = value
		case "BASIC_AUTH_PASSWORD":
			tunnel.BasicAuth.Password = value
		case "BASIC_AUTH_PASSWORD_ENV":
			tunnel.BasicAuth.PasswordEnv = value
		default:
			return nil, fmt.Errorf("unsupported tunnel env field %q", key)
		}
	}

	if len(indexed) == 0 {
		return nil, nil
	}

	indices := make([]int, 0, len(indexed))
	for idx := range indexed {
		indices = append(indices, idx)
	}
	sort.Ints(indices)

	tunnels := make([]TunnelConfig, 0, len(indices))
	for _, idx := range indices {
		tunnel := *indexed[idx]
		if tunnel.Name == "" && tunnel.LocalAddr == "" && tunnel.Hostname == "" && tunnel.Access == "" {
			continue
		}
		if tunnel.Name == "" || tunnel.LocalAddr == "" {
			return nil, fmt.Errorf("BLOOP_TUNNELS_%d requires NAME and LOCAL_ADDR", idx)
		}
		tunnels = append(tunnels, tunnel)
	}
	return tunnels, nil
}
