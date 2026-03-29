# Data Model: BubbleTea TUI for bloop-client Setup

**Feature**: [003-bubbletea-tui-for-client-setup]  
**Date**: 2026-03-29  
**Purpose**: Define data structures for BubbleTea models, screen states, and configuration

---

## Overview

This document defines the data structures used throughout the BubbleTea TUI implementation. These structures represent the pure model state that is updated by the Update function and rendered by the View function, following the Elm architecture pattern.

---

## Core Types

### State Enum

Represents the current screen in the TUI navigation flow.

```go
type State int

const (
    StateWelcome State = iota
    StateConfig
    StateTunnels
    StateVerification
    StateDockerDiscovery
    StateReview
    StateOutput
    StateQuitConfirm
)
```

**Transitions**:
- Welcome → Config
- Config → Tunnels (or Verification if verification already run)
- Tunnels → Docker Discovery (optional) → Verification → Review
- Verification → Review (or back to Config)
- Review → Output (or back to any previous screen)
- Any screen → QuitConfirm (if unsaved changes)
- QuitConfirm → Exit or return to previous screen

### Config Structure

Represents the complete client configuration being edited.

```go
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
    LocalAddr  string
    Access     string // "public", "basic_auth", "token_protected"
    BasicAuth  BasicAuthConfig
    TokenEnv   string // For token_protected access
}

type BasicAuthConfig struct {
    Username    string
    PasswordEnv string // Environment variable name
}
```

**Validation Rules**:
- `ControlPlaneURL`: Required, must be valid URL
- `RelayURL`: Required, must be valid URL
- `AuthTokenEnv`: Required (or AuthToken in non-TUI mode)
- `EnrollmentToken` OR `EnrollmentTokenEnv`: Optional, at least one recommended
- `Reconnect.InitialDelayMs`: Required, must be > 0
- `Reconnect.MaxDelayMs`: Required, must be > InitialDelayMs
- `Tunnels`: Required, at least one tunnel
- `Tunnels[i].Name`: Required, non-empty
- `Tunnels[i].LocalAddr`: Required, non-empty
- `Tunnels[i].Access`: Required, one of: "public", "basic_auth", "token_protected"
- `Tunnels[i].BasicAuth`: Required if Access = "basic_auth"
- `Tunnels[i].TokenEnv`: Required if Access = "token_protected"

### OutputMode Enum

Represents the selected output format for generated configuration.

```go
type OutputMode string

const (
    OutputYAML         OutputMode = "yaml"
    OutputEnvFile      OutputMode = "env-file"
    OutputComposeBlock OutputMode = "compose-block"
)
```

---

## BubbleTea Models

### MainModel

Root model that manages screen state and navigation stack.

```go
type MainModel struct {
    state          State
    stateStack     []State          // For back navigation
    config         Config
    currentScreen  tea.Model        // Current screen model
    previousScreen tea.Model        // Previous screen (for back navigation)
    isLoading      bool
    errorMsg       string
    width          int              // Terminal width
    height         int              // Terminal height
}

type ScreenTransitionMsg struct {
    from, to State
}

type ConfigChangedMsg struct{}
```

**State Transitions**:
- `PushScreen(to State)`: Add current state to stack, transition to new state
- `PopScreen()`: Pop previous state from stack, transition to previous state
- `ResetScreens()`: Clear stack, reset to StateWelcome

### InputFieldModel

Reusable text/password input field model.

```go
type InputFieldModel struct {
    label         string
    placeholder   string
    value         string
    focused       bool
    isPassword    bool
    validation    func(string) error
    errorMsg      string
    cursorPos     int
}

type InputFieldFocusedMsg struct{ Index int }
type InputFieldBlurMsg struct{ Index int }
type InputFieldValueChangeMsg struct{ Index int; Value string }
```

**Validation**:
- Optional validation function called on each value change
- `errorMsg` displayed below field if validation fails

### SelectFieldModel

Reusable selection field model (dropdown/radio buttons).

```go
type SelectFieldModel struct {
    label         string
    options       []string
    selectedIndex int
    focused       bool
}

type SelectFieldChangeMsg struct{ Index int }
```

### ListViewModel

Reusable scrollable list with selection.

```go
type ListViewModel struct {
    items         []ListItem
    selectedIndex  int
    cursorPos     int    // Visible scrolling position
    focused       bool
    height        int    // Number of visible items
}

type ListItem struct {
    ID       string
    Label    string
    Details  string // Optional additional info
}

type ListSelectMsg struct{ Index int }
type ListDeleteMsg struct{ Index int }
type ListEditMsg struct{ Index int }
```

**Scrolling Logic**:
- When `selectedIndex` moves beyond visible window, scroll window down
- When `selectedIndex` moves before visible window, scroll window up
- Window size determined by `height`

### StatusModel

Loading/spinner indicator for async operations.

```go
type StatusModel struct {
    loading     bool
    message     string
    spinner     spinner.Model
}

type StatusLoadingMsg struct{ Message string }
type StatusSuccessMsg struct{ Message string }
type StatusErrorMsg struct{ Message string }
```

---

## Screen Models

### WelcomeModel

Welcome screen with quick start options.

```go
type WelcomeModel struct {
    options      []WelcomeOption
    selectedIndex int
    loadedConfig Config // If editing existing
}

type WelcomeOption struct {
    ID          string
    Label       string
    Description string
    Action      func() tea.Cmd
}

type WelcomeChoiceMsg struct{ Choice string }
```

**Options**:
1. "New Setup" - Clear defaults
2. "Edit Existing" - Load from config path
3. "Quick Start" - Skip to output with defaults

### ConfigModel

Configuration form screen.

```go
type ConfigModel struct {
    controlPlaneURL    InputFieldModel
    relayURL           InputFieldModel
    authTokenEnv        InputFieldModel
    enrollmentToken    InputFieldModel
    enrollmentTokenEnv InputFieldModel
    reconnectInitDelay InputFieldModel
    reconnectMaxDelay  InputFieldModel
    testConnectivity   bool
    connectivityResult ConnectivityResult
    errors             map[string]string
}

type ConfigSaveMsg struct{ Config Config }
type ConfigTestConnectivityMsg struct{}
```

**Field Focus Order**:
1. Control Plane URL
2. Relay URL
3. Auth Token Env
4. Enrollment Token
5. Enrollment Token Env
6. Reconnect Initial Delay
7. Reconnect Max Delay

### TunnelsListModel

List view for managing tunnels.

```go
type TunnelsListModel struct {
    tunnels      []TunnelConfig
    listModel    ListViewModel
    errorMessage string
}

type TunnelsListAddMsg struct{}
type TunnelsListEditMsg struct{ Index int }
type TunnelsListDeleteMsg struct{ Index int }
type TunnelsListBackMsg struct{}
type TunnelsListNextMsg struct{}
```

**List Columns**:
- Name
- Local Addr
- Access
- Hostname (if any)

### TunnelFormModel

Form for editing a single tunnel.

```go
type TunnelFormModel struct {
    tunnel       TunnelConfig
    isEditing    bool // true if editing existing, false if new
    name         InputFieldModel
    hostname     InputFieldModel
    localAddr    InputFieldModel
    accessType   SelectFieldModel
    basicAuthUsername   InputFieldModel
    basicAuthPasswordEnv InputFieldModel
    tokenEnv     InputFieldModel
    errors       map[string]string
}

type TunnelFormSaveMsg struct{ Tunnel TunnelConfig }
type TunnelFormCancelMsg struct{}
```

**Conditional Fields**:
- `basicAuthUsername`, `basicAuthPasswordEnv`: Visible only when access = "basic_auth"
- `tokenEnv`: Visible only when access = "token_protected"

### DockerDiscoveryModel

Screen for discovering Docker services.

```go
type DockerDiscoveryModel struct {
    candidates   []DockerCandidate
    selected     map[int]bool // Map of candidate index to selected state
    listModel    ListViewModel
    isLoading    bool
    errorMessage string
    discovered  bool
}

type DockerCandidate struct {
    ContainerID   string
    ContainerName string
    ImageName     string
    CandidatePorts []PortMapping
}

type PortMapping struct {
    ContainerPort int
    HostPort      int
    Protocol      string // "tcp", "udp"
    IsPublished   bool
}

type DockerDiscoveryMsg struct{ Candidates []DockerCandidate }
type DockerDiscoverySelectMsg struct{ Index int }
type DockerDiscoveryAddMsg struct{}
type DockerDiscoverySkipMsg struct{}
```

**Candidate Selection**:
- Multiple candidates can be selected
- Selected candidates will be added as tunnels to main config
- Each candidate becomes a tunnel with:
  - Name: `{ContainerName}-{Port}`
  - LocalAddr: `{HostIP}:{HostPort}` or `127.0.0.1:{ContainerPort}`
  - Access: "public"

### VerificationModel

Screen for running API verification tests.

```go
type VerificationModel struct {
    tabs         []VerificationTab
    selectedTab  int
    results      VerificationResults
    isRunning    bool
}

type VerificationTab struct {
    ID      string
    Label   string
    Status  VerificationStatus
}

type VerificationStatus string

const (
    VerificationStatusPending   VerificationStatus = "pending"
    VerificationStatusRunning   VerificationStatus = "running"
    VerificationStatusSuccess   VerificationStatus = "success"
    VerificationStatusError     VerificationStatus = "error"
    VerificationStatusSkipped   VerificationStatus = "skipped"
)

type VerificationResults struct {
    Connectivity ConnectivityResult
    Enrollment  EnrollmentResult
    Relay       RelayResult
}

type ConnectivityResult struct {
    Status    VerificationStatus
    Latency   time.Duration
    DNSError  error
    TLSError  error
    HTTPError error
    Timestamp time.Time
}

type EnrollmentResult struct {
    Status         VerificationStatus
    InstallationID string
    IngestToken    string // Not displayed
    Error          error
    ErrorCode      string // "invalid_token", "expired", "unauthorized"
    Timestamp      time.Time
}

type RelayResult struct {
    Status    VerificationStatus
    Error     error
    Details   string // WebSocket close code, timeout, etc.
    Timestamp time.Time
}

type VerificationRunMsg struct{}           // Trigger all verifications
type VerificationRetryMsg struct{ TabID string }
type VerificationSkipMsg struct{}
type VerificationNextMsg struct{}
```

**Tabs**:
1. "Control Plane Connectivity"
2. "Enrollment Token"
3. "Relay Auth Token"

### ReviewModel

Summary screen showing all configuration.

```go
type ReviewModel struct {
    config         Config
    sections       []ReviewSection
}

type ReviewSection struct {
    Title    string
    Content  string // Formatted summary
    EditCmd  tea.Cmd // Command to navigate to edit screen
}

type ReviewEditMsg struct{ Section string }
type ReviewFinishMsg struct{}
type ReviewBackMsg struct{}
```

**Sections**:
1. Control Plane (URL, verification status)
2. Relay (URL, token env, verification status)
3. Enrollment (token status, installation ID if verified)
4. Tunnels (list of all tunnels)
5. Output (mode, path)

### OutputModel

Screen for selecting output mode and generating output.

```go
type OutputModel struct {
    outputMode   OutputMode
    outputPath   string
    outputPathFocused bool
    isWriting    bool
    errorMessage string
    output       string // Generated output for display
    outputShown  bool
}

type OutputSelectModeMsg struct{ Mode OutputMode }
type OutputPathChangeMsg struct{ Path string }
type OutputGenerateMsg struct{}
type OutputFinishMsg struct{}
type OutputBackMsg struct{}
```

### QuitConfirmModel

Confirmation modal for exiting with unsaved changes.

```go
type QuitConfirmModel struct {
    message string
    options []string
    selected int
}

type QuitConfirmMsg struct{ Choice string }
```

**Options**:
- "Quit"
- "Cancel"

---

## Messages

### BubbleTea Standard Messages

- `tea.KeyMsg`: Keyboard input
- `tea.MouseMsg`: Mouse input (if supported)
- `tea.WindowSizeMsg`: Terminal resize
- `tea.QuitMsg`: Quit command

### Custom Messages

```go
// Screen navigation
type ScreenTransitionMsg struct {
    from, to State
}

// Config updates
type ConfigChangedMsg struct{}
type ConfigSaveMsg struct{ Config Config }

// Input field updates
type InputFieldFocusedMsg struct{ Index int }
type InputFieldBlurMsg struct{ Index int }
type InputFieldValueChangeMsg struct{ Index int; Value string }

// List actions
type ListSelectMsg struct{ Index int }
type ListDeleteMsg struct{ Index int }
type ListEditMsg struct{ Index int }

// Status updates
type StatusLoadingMsg struct{ Message string }
type StatusSuccessMsg struct{ Message string }
type StatusErrorMsg struct{ Message string }

// API verification
type VerificationRunMsg struct{}
type VerificationRetryMsg struct{ TabID string }
type VerificationSkipMsg struct{}

// Docker discovery
type DockerDiscoverySelectMsg struct{ Index int }
type DockerDiscoveryAddMsg struct{}
type DockerDiscoverySkipMsg struct{}

// Output
type OutputSelectModeMsg struct{ Mode OutputMode }
type OutputGenerateMsg struct{}

// Async operations
type TickMsg time.Time
```

---

## Validation

### Validation Functions

```go
// URL validation
func ValidateURL(value string) error {
    u, err := url.Parse(value)
    if err != nil {
        return fmt.Errorf("invalid URL: %w", err)
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return fmt.Errorf("URL must use http or https scheme")
    }
    if u.Host == "" {
        return fmt.Errorf("URL must include a host")
    }
    return nil
}

// Port validation
func ValidatePort(value string) error {
    port, err := strconv.Atoi(value)
    if err != nil {
        return fmt.Errorf("must be a number")
    }
    if port < 1 || port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    return nil
}

// Positive integer validation
func ValidatePositiveInt(value string) error {
    num, err := strconv.Atoi(value)
    if err != nil {
        return fmt.Errorf("must be a number")
    }
    if num <= 0 {
        return fmt.Errorf("must be greater than 0")
    }
    return nil
}

// Required field validation
func ValidateRequired(value string, fieldName string) error {
    if strings.TrimSpace(value) == "" {
        return fmt.Errorf("%s is required", fieldName)
    }
    return nil
}
```

---

## Security Considerations

### Sensitive Data Handling

**Auth Tokens**:
- Never store auth tokens in plaintext in model state
- Use password field type for token inputs (masked)
- Clear token values from memory after use
- Redact tokens from logs and error messages

**Config File Permissions**:
- Write config files with 0o600 permissions
- Verify permissions after writing
- Fail if permissions cannot be set

**API Requests**:
- Validate URLs before making requests (prevent SSRF)
- Use appropriate timeouts (5s connectivity, 10s enrollment/relay)
- Redact auth tokens from request headers in logs

### In-Memory Cleanup

```go
// Clear sensitive values from memory
func ClearSensitiveData(cfg *Config) {
    cfg.EnrollmentToken = ""
    cfg.IngestToken = ""
    for i := range cfg.Tunnels {
        if cfg.Tunnels[i].Access == "token_protected" {
            // Clear token if stored in plaintext (should use env var)
        }
    }
}
```

---

## State Persistence

### Config Loading

Config is loaded from existing file when "Edit Existing" is selected:

```go
func LoadConfig(path string) (Config, error) {
    cfg, err := config.LoadClientConfig(path)
    if err != nil {
        return Config{}, err
    }

    return Config{
        ControlPlaneURL:      cfg.ControlPlaneURL,
        RelayURL:             cfg.RelayURL,
        AuthTokenEnv:         cfg.AuthTokenEnv,
        EnrollmentTokenEnv:   cfg.EnrollmentTokenEnv,
        Reconnect: ReconnectConfig{
            InitialDelayMs: cfg.Reconnect.InitialDelayMs,
            MaxDelayMs:     cfg.Reconnect.MaxDelayMs,
        },
        Tunnels:              cfg.Tunnels,
        OutputMode:           OutputYAML,
        OutputPath:           path,
        HasChanges:           false,
    }, nil
}
```

### Output Generation

Config is serialized to selected output format:

```go
func GenerateOutput(cfg Config, mode OutputMode, outputPath string) (string, error) {
    switch mode {
    case OutputYAML:
        return generateYAML(cfg, outputPath)
    case OutputEnvFile:
        return generateEnvFile(cfg, false)
    case OutputComposeBlock:
        return generateEnvFile(cfg, true)
    default:
        return "", fmt.Errorf("unsupported output mode: %s", mode)
    }
}
```

---

## Data Flow

### Configuration Flow

```
User Input
    ↓
InputFieldModel Update
    ↓
Config Update (ConfigChangedMsg)
    ↓
Validation
    ↓
Error Display (if invalid)
    ↓
State Update (if valid)
```

### API Verification Flow

```
User clicks "Verify All"
    ↓
VerificationModel Update
    ↓
Spawn Async Commands (goroutines)
    ↓
Status Updates (StatusLoadingMsg)
    ↓
API Requests
    ↓
Results (APITestMsg)
    ↓
VerificationResults Update
    ↓
Status Display (Success/Error)
```

### Screen Navigation Flow

```
User Action (e.g., Enter on "Next")
    ↓
Screen Model Update
    ↓
ScreenTransitionMsg
    ↓
MainModel Update
    ↓
Push/Pop State Stack
    ↓
Transition to New Screen
    ↓
New Screen Model Init
    ↓
View Render
```

---

## Notes

- All models must be pure (no side effects in model code)
- Side effects happen via `tea.Cmd` returned from Update functions
- Views must be deterministic (same model → same output)
- Use `TickMsg` for periodic status updates during async operations
- Screen stack enables back navigation without losing intermediate state
- `HasChanges` flag tracks whether user has modified config since load
