package tui

// Screen state enum
type State int

const (
	StateWelcome State = iota
	StateConfig
	StateEndpoints
	StateTunnels
	StateVerification
	StateDockerDiscovery
	StateReview
	StateOutput
	StateQuitConfirm
)

// Shared configuration structure
type Config struct {
	ControlPlaneURL      string
	RelayURL             string
	AuthTokenEnv         string
	EnrollmentToken      string // Optional
	EnrollmentTokenEnv   string // Optional
	Reconnect            ReconnectConfig
	Tunnels              []TunnelConfig
	OutputMode           OutputMode
	OutputPath           string
	HasChanges           bool // Track unsaved changes
	VerificationResults  VerificationResults
	InstallationID       string // From enrollment verification
	IngestToken          string // From enrollment verification (not displayed)
}

type ReconnectConfig struct {
	InitialDelayMs int
	MaxDelayMs     int
}

type TunnelConfig struct {
	Name       string
	Hostname   string // Optional
	LocalIP    string // IP address or hostname
	LocalPort  int    // Port number (1-65535)
	Access     string // "public", "basic_auth", "token_protected"
	BasicAuth  BasicAuthConfig
	TokenEnv   string // For token_protected access
}

type BasicAuthConfig struct {
	Username    string
	PasswordEnv string // Environment variable name
}

type VerificationResults struct {
	Connectivity ConnectivityResult
	Enrollment  EnrollmentResult
	Relay       RelayResult
}

type ConnectivityResult struct {
	Success    bool
	Status     string // "pending", "running", "success", "error", "skipped"
	Latency    string // Formatted duration
	DNSError  error
	TLSError  error
	HTTPError  error
}

type EnrollmentResult struct {
	Success        bool
	Status         string // "pending", "running", "success", "error", "skipped"
	InstallationID string
	IngestToken    string // Not displayed
	Error          error
	ErrorCode      string // "invalid_token", "expired", "unauthorized"
}

type RelayResult struct {
	Success  bool
	Status   string // "pending", "running", "success", "error", "skipped"
	Error     error
	Details  string // WebSocket close code, timeout, etc.
}

type OutputMode string

const (
	OutputYAML         OutputMode = "yaml"
	OutputEnvFile      OutputMode = "env-file"
	OutputComposeBlock OutputMode = "compose-block"
)

// Messages

type ScreenTransitionMsg struct {
	From, To State
}

type ConfigChangedMsg struct{}

type ConfigSaveMsg struct{ Config }

type TunnelAddMsg struct{}
type TunnelEditMsg struct{ Index int }
type TunnelDeleteMsg struct{ Index int }
type TunnelSaveMsg struct{ Tunnel TunnelConfig }

type TunnelsListAddMsg struct{}
type TunnelsListEditMsg struct{ Index int }
type TunnelsListDeleteMsg struct{ Index int }
type TunnelsListBackMsg struct{}
type TunnelsListNextMsg struct{}

type QuitConfirmMsg struct{ Choice string }
