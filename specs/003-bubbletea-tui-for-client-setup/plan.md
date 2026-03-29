# Implementation Plan: BubbleTea TUI for bloop-client Setup

**Branch**: `[feat/bloop-client-bubbletea-tui]` | **Date**: 2026-03-29 | **Spec**: `/specs/003-bubbletea-tui-for-client-setup/spec.md`
**Input**: Feature specification from `/specs/003-bubbletea-tui-for-client-setup/spec.md`

## Summary

Upgrade the existing bufio-based `bloop-tunnel setup` command to use BubbleTea framework for a modern, responsive terminal UI. The TUI will provide screen-based navigation, real-time form validation, API verification flows (connectivity, enrollment, relay), and secure token handling. All existing functionality (config loading, output modes, non-interactive mode) will be preserved while providing a significantly improved operator experience.

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**: BubbleTea (github.com/charmbracelet/bubbletea), Lipgloss (github.com/charmbracelet/lipgloss), existing config/auth packages
**Storage**: File-based client config (YAML) and environment variables; no database changes
**Testing**: Go test with table-driven tests for model state transitions, integration tests for full setup flows
**Target Platform**: Native macOS/Linux/Windows CLI terminals; Docker container runtime
**Project Type**: CLI application upgrade from linear prompts to screen-based TUI
**Performance Goals**: 60fps rendering, API verification completes in <10s, TUI responds to input within 50ms
**Constraints**: Backward compatible with existing flags and config format; must degrade to non-interactive mode when TTY unavailable
**Scale/Scope**: Single-operator setup experience; unlimited tunnel definitions supported

## Constitution Check

### Initial Gate Review

- **Operator-First Design**: Pass. Screen-based TUI with intuitive navigation provides better UX than linear prompts.
- **Production-Ready Defaults**: Pass. All defaults continue to point to api.bloop.to; TUI defaults preserve existing behavior.
- **Test-First (NON-NEGOTIABLE)**: Pass. All BubbleTea models will have table-driven tests for state transitions.
- **Security by Default**: Pass. Auth tokens masked in input, not displayed in plaintext, config files written with 0o600 permissions.
- **Backward Compatibility**: Pass. All existing flags (--non-interactive, --output, --config) preserved.
- **Minimal Dependencies**: Pass. BubbleTea and Lipgloss are well-maintained, no additional dependencies added.
- **Graceful Degradation**: Pass. TUI falls back to non-interactive mode when TTY unavailable, all errors handled gracefully.

No constitution violations currently require justification.

## Project Structure

### Documentation (this feature)

```text
specs/003-bubbletea-tui-for-client-setup/
├── plan.md              # This file
├── spec.md              # Feature specification (already created)
├── tasks.md             # Implementation tasks
├── data-model.md        # BubbleTea model state definitions
└── contracts/           # API contracts for verification endpoints
    ├── control-plane.md
    └── relay.md
```

### Source Code (repository root)

```text
internal/client/tui/              # New BubbleTea TUI package
├── main.go                       # TUI entry point, model initialization
├── styles.go                     # Lipgloss style definitions
├── types.go                      # Shared types and messages
├── screens/                      # Screen-specific models and views
│   ├── welcome.go
│   ├── config.go
│   ├── tunnels.go
│   ├── tunnel_form.go
│   ├── docker_discovery.go
│   ├── verification.go
│   ├── review.go
│   └── output.go
├── models/                       # Reusable models
│   ├── input_field.go            # Text/password input field
│   ├── select_field.go           # Selection field (dropdown/radio)
│   ├── list_view.go              # Scrollable list view
│   └── status.go                 # Loading/progress indicator
├── api/                          # API verification logic
│   ├── connectivity.go           # Control plane connectivity test
│   ├── enrollment.go             # Enrollment token verification
│   └── relay.go                  # Relay auth token verification
└── util/                         # TUI utilities
    ├── keys.go                   # Keyboard event handling
    └── clipboard.go              # Clipboard operations (optional)

internal/client/setup/            # Existing setup package (modified)
├── setup.go                      # Will call TUI or fall back to prompts
└── output.go                     # Reused for YAML/env/compose output

cmd/bloop-client/
└── main.go                       # Modified to route setup to TUI

tests/
├── tui/                          # TUI-specific tests
│   ├── models_test.go            # Model state transition tests
│   ├── screens_test.go           # Screen-specific tests
│   └── integration_test.go       # Full setup flow tests
└── integration/                  # Existing integration tests (preserved)
```

**Structure Decision**: Create a new `internal/client/tui` package for the BubbleTea implementation to keep it isolated from the existing setup code. This allows for gradual migration and preserves the non-interactive fallback. The `internal/client/setup` package will be modified to route to the TUI when interactive mode is requested and TTY is available.

## Architecture Overview

### BubbleTea Model-View-Update Architecture

**Model**: Pure data structures representing TUI state. No side effects, no I/O.

**View**: Functions that render the model to terminal output using Lipgloss styling. Deterministic, no side effects.

**Update**: Functions that handle incoming messages (tea.Msg) and return updated models and optional commands. This is where all side effects (API calls, file I/O) happen via tea.Cmd.

### Screen Navigation Flow

```
┌─────────────┐
│   Welcome   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Config    │◄─────────────────┐
└──────┬──────┘                  │
       │                         │
       ▼                         │
┌─────────────┐                  │
│   Tunnels   │                  │
└──────┬──────┘                  │
       │                         │
       ▼                         │
┌─────────────┐                  │
│   Verify    │◄─────────────────┤
│    APIs     │                  │
└──────┬──────┘                  │
       │                         │
       ▼                         │
┌─────────────┐                  │
│   Docker    │ (optional)       │
│  Discovery  │───────────────────┤
└──────┬──────┘                  │
       │                         │
       ▼                         │
┌─────────────┐                  │
│   Review    │───────────────────┘
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Output    │
└─────────────┘
```

**Navigation patterns**:
- **Forward**: Enter/Space on confirmation buttons, auto-advance after successful operations
- **Back**: Esc on non-modal screens, back buttons in review/verification screens
- **Jump**: Tab to skip to next section (optional enhancement)
- **Quit**: Q or Ctrl+C triggers confirmation modal if unsaved changes

### State Management

**MainModel State**:
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

**Config State**:
```go
type Config struct {
    ControlPlaneURL   string
    RelayURL          string
    AuthTokenEnv      string
    EnrollmentToken   string
    EnrollmentTokenEnv string
    Reconnect         ReconnectConfig
    Tunnels           []TunnelConfig
    OutputMode        OutputMode
    OutputPath        string
    HasChanges        bool
}

type ReconnectConfig struct {
    InitialDelayMs int
    MaxDelayMs     int
}
```

**Screen State Stack**:
Each screen maintains its own state. Navigation pushes/pops screens onto a stack, allowing back navigation without losing intermediate state.

### API Verification Design

**Connectivity Test**:
```go
type ConnectivityResult struct {
    Success      bool
    Latency      time.Duration
    DNSError     error
    TLSError     error
    HTTPError    error
}
```

Flow:
1. Resolve DNS hostname
2. Perform TLS handshake
3. Send HTTP GET to `/health` or root
4. Report result with specific error details

**Enrollment Token Verification**:
```go
type EnrollmentResult struct {
    Success        bool
    InstallationID string
    IngestToken    string
    Error          error
    ErrorCode      string // "invalid_token", "expired", "unauthorized"
}
```

Flow:
1. Send POST to `{control_plane}/api/runtime/enroll`
2. Parse response
3. Return installation ID and ingest token on success
4. Return specific error code on failure

**Relay Token Verification**:
```go
type RelayResult struct {
    Success bool
    Error   error
    Details string // WebSocket close code, timeout, etc.
}
```

Flow:
1. Establish WebSocket connection to relay URL with auth token
2. Send minimal handshake message
3. Verify authentication success/failure
4. Close connection

### Screen Designs

**Welcome Screen**:
- Title: "bloop-client Setup"
- Subtitle: "Configure your tunnel client"
- Options:
  - "New Setup" (clear defaults)
  - "Edit Existing" (load from config path)
  - "Quick Start" (skip to output with defaults)
- Footer: "Press q to quit, arrow keys to navigate"

**Config Screen**:
- Fields:
  - Control Plane URL (text input, default: https://api.bloop.to)
  - Relay URL (text input, default: production relay)
  - Relay Auth Token Env (text input, default: BLOOP_CLIENT_TOKEN)
  - Enrollment Token (password input, optional)
  - Enrollment Token Env (text input, optional)
  - Reconnect Initial Delay (number input, default: 1000)
  - Reconnect Max Delay (number input, default: 30000)
- Validation:
  - URLs must be valid and non-empty
  - Delay values must be positive integers
  - At least one of token or token env must be specified
- Actions:
  - "Next" → Tunnels screen
  - "Back" → Welcome screen
  - "Test Connectivity" → Run connectivity test, show inline result

**Tunnels List Screen**:
- List view showing all configured tunnels
- Columns: Name, Local Addr, Access, Hostname (if any)
- Actions:
  - "Add Tunnel" → Tunnel Form screen (new tunnel)
  - "Edit" → Tunnel Form screen (existing tunnel)
  - "Remove" → Remove tunnel with confirmation
  - "Next" → Verification screen
  - "Back" → Config screen

**Tunnel Form Screen**:
- Fields:
  - Name (text input, required)
  - Local IP Address (text input, required)
  - Port Number (number input, required, 1-65535)
  - Hostname (text input, optional)
  - Access Type (select: public, basic_auth, token_protected)
  - Basic Auth Username (text input, conditional on access=basic_auth)
  - Basic Auth Password Env (text input, conditional on access=basic_auth)
  - Token Env (text input, conditional on access=token_protected)
- Validation:
  - Name and Local Addr required
  - Conditional fields required based on access type
- Actions:
  - "Save" → Return to Tunnels List
  - "Cancel" → Return to Tunnels List without saving

**Docker Discovery Screen** (optional):
- List view of discovered containers with candidate ports
- Columns: Container Name, Image, Candidate Ports, Actions
- Actions:
  - "Select" → Mark as tunnel to add
  - "Deselect" → Unmark
  - "Add Selected" → Add marked containers as tunnels to main config
  - "Skip" → Continue without adding
  - "Refresh" → Re-run discovery
- Discovery error handling:
  - Show clear error if Docker socket unavailable
  - Offer "Retry" and "Skip" options
  - Graceful degradation if Docker not running

**Verification Screen**:
- Tabs or sections for each verification:
  1. Control Plane Connectivity
  2. Enrollment Token
  3. Relay Auth Token
- Each section shows:
  - Status indicator (pending/running/success/error)
  - Progress spinner or checkmark
  - Error details if failed
  - "Retry" button for failed verifications
- Actions:
  - "Verify All" → Run all verifications in sequence
  - "Skip" → Skip verification and continue
  - "Next" → Continue to Review (may warn if verifications skipped/failed)
  - "Back" → Return to previous screen

**Review Screen**:
- Summary sections:
  - Control Plane: URL, verified status
  - Relay: URL, token env, verified status
  - Enrollment: token status (if configured), installation ID (if verified)
  - Tunnels: List of all tunnels with details
  - Output: Mode, path (if file output)
- Actions:
  - "Edit Config" → Return to Config screen
  - "Edit Tunnels" → Return to Tunnels List screen
  - "Change Output" → Return to Output screen
  - "Finish & Generate" → Generate output and exit

**Output Screen**:
- Options:
  - "Write YAML to file" → Prompt for file path, write config
  - "Print env-file to stdout" → Print .env format
  - "Print Compose block to stdout" → Print Docker Compose env block
- For file output:
  - Prompt for path with default from flag
  - Write with 0o600 permissions
  - Show success message with path
- For stdout output:
  - Print to terminal
  - Show "Output written above" message
  - Offer option to save to file
- Actions:
  - "Finish" → Exit TUI
  - "Back" → Return to Review screen

**Quit Confirmation Modal**:
- Message: "You have unsaved changes. Are you sure you want to quit?"
- Options:
  - "Quit" → Exit TUI without saving
  - "Cancel" → Return to current screen

### Styling Design (Lipgloss)

**Color Palette** (accessible, works in light/dark):
- Primary: Blue (#007AFF)
- Success: Green (#34C759)
- Error: Red (#FF3B30)
- Warning: Orange (#FF9500)
- Muted: Gray (#8E8E93)
- Background: Terminal default (transparent)

**Component Styles**:
```go
var (
    titleStyle       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
    subtitleStyle    = lipgloss.NewStyle().Foreground(mutedColor)
    inputStyle       = lipgloss.NewStyle().Foreground(textColor)
    inputFocusedStyle = inputStyle.Border(lipgloss.RoundedBorder())
    buttonStyle      = lipgloss.NewStyle().Padding(0, 2)
    buttonFocusedStyle = buttonStyle.Foreground(white).Background(primaryColor)
    successStyle     = lipgloss.NewStyle().Foreground(successColor)
    errorStyle       = lipgloss.NewStyle().Foreground(errorColor)
    listItemStyle    = lipgloss.NewStyle().Padding(0, 1)
    listFocusedStyle = listItemStyle.Background(mutedColor).Foreground(white)
)
```

**Layout**:
- Centered content with 10-20 character padding
- Max width of 120 characters for readability
- Sections separated by blank lines or dividers
- Footer with keyboard shortcuts consistently positioned

## Integration with Existing Codebase

### Config Loading

Reuse existing `config.LoadClientConfig()` and `config.DefaultClientConfig()`. The TUI will:
1. Check for existing config at path from flag
2. Load and merge with defaults
3. Preload values into form fields for editing

### Output Generation

Reuse existing output logic from `setup.output.go`:
- `renderEnv()` for env-file and compose-block formats
- YAML marshaling for file output
- Validation logic from `validateConfig()`

### Docker Discovery

Reuse existing `dockerdiscover` package:
- `Discoverer` interface for candidate enumeration
- Candidate filtering logic
- Error handling for unavailable socket

### Non-Interactive Mode

Preserve existing non-interactive behavior:
- Detect TTY availability
- If no TTY or --non-interactive flag, use existing prompt-based setup
- TUI only runs when TTY available and interactive mode requested

## Testing Strategy

### Unit Tests

**Model State Transitions**:
```go
func TestConfigModelUpdate(t *testing.T) {
    tests := []struct {
        name     string
        initial  ConfigModel
        msg      tea.Msg
        expected ConfigModel
        cmd      tea.Cmd
    }{
        // Test cases for each input field update
        // Test validation logic
        // Test error handling
    }
    // Table-driven test implementation
}
```

**Screen Navigation**:
- Test forward and backward navigation
- Test screen stack push/pop
- Test quit confirmation with/without changes

**Input Field Behavior**:
- Test text input updates
- Test password masking
- Test validation feedback

### Integration Tests

**Full Setup Flow**:
- Start TUI in headless mode
- Simulate keyboard input for complete setup
- Verify output matches expected config

**API Verification**:
- Mock control plane API responses
- Mock relay WebSocket connections
- Test success and error paths

**Docker Discovery**:
- Mock Docker client with test containers
- Test candidate selection and filtering
- Test error handling for unavailable socket

### Manual Testing

**Terminal Compatibility**:
- Test on macOS Terminal, iTerm2
- Test on Linux gnome-terminal, xterm
- Test on Windows Terminal, PowerShell
- Test with different terminal sizes (80x24, 120x40)

**Accessibility**:
- Test with screen readers (NVDA, VoiceOver)
- Verify all interactive elements have text labels
- Verify keyboard-only navigation works

**Color Schemes**:
- Test with light terminal background
- Test with dark terminal background
- Test with high-contrast settings

## Performance Considerations

**Rendering**:
- Use efficient diff-based rendering (BubbleTea handles this)
- Minimize re-renders by updating model only on actual changes
- Lazy-render lists for large tunnel counts

**Async Operations**:
- API verifications run in background goroutines
- Use TickMsg for status updates
- Set appropriate timeouts (5s connectivity, 10s enrollment/relay)

**Memory**:
- Keep models lightweight (no large data structures)
- Clear sensitive values from memory after use
- Reuse string builders for output generation

## Security Considerations

**Token Handling**:
- Mask all token inputs (password field type)
- Never display tokens in plaintext in TUI
- Clear token values from memory after use

**Config File Permissions**:
- Write config files with 0o600 permissions
- Verify permissions after writing
- Fail if permissions cannot be set

**API Verification**:
- Use appropriate timeouts to prevent hanging
- Validate URLs before making requests (prevent SSRF)
- Redact tokens from logs and error messages

**Docker Discovery**:
- Always require explicit confirmation before adding tunnels
- Never automatically expose discovered services
- Show clear warning before Docker socket access

## Implementation Phases

### Phase 1 - Foundations
- Add BubbleTea and Lipgloss dependencies
- Create `internal/client/tui` package structure
- Implement basic TUI entry point with empty MainModel
- Implement styling system (styles.go)
- Implement shared types and messages (types.go)
- Add basic reusable models (input_field, select_field)

### Phase 2 - Core Screens
- Implement Welcome screen
- Implement Config screen with all fields
- Implement Config screen validation
- Implement Tunnels List screen
- Implement Tunnel Form screen
- Implement screen navigation between core screens

### Phase 3 - API Verification
- Implement control plane connectivity test
- Implement enrollment token verification
- Implement relay auth token verification
- Implement Verification screen combining all tests
- Add error handling and retry logic

### Phase 4 - Docker Discovery
- Integrate existing dockerdiscover package
- Implement Docker Discovery screen
- Implement candidate selection and confirmation
- Add error handling for unavailable socket

### Phase 5 - Review and Output
- Implement Review screen with summary
- Implement Output screen with mode selection
- Implement file output with 0o600 permissions
- Implement stdout output for env/compose formats
- Implement quit confirmation modal

### Phase 6 - Integration and Polish
- Integrate TUI with existing setup command
- Implement TTY detection and fallback
- Preserve all existing flags
- Add keyboard shortcuts
- Implement terminal resize handling
- Add loading states and spinners

### Phase 7 - Testing
- Write unit tests for all models
- Write integration tests for full flows
- Write tests for API verification
- Write tests for Docker discovery
- Manual testing on multiple terminals

### Phase 8 - Security Review
- Spawn security review subagent
- Audit token handling
- Audit API verification security
- Audit file permissions handling
- Implement any findings from review

### Phase 9 - Documentation
- Update README with TUI screenshots
- Write TUI usage guide
- Document keyboard shortcuts
- Add troubleshooting section
- Update CLIENT_INSTALL.md

### Phase 10 - Merge
- Verify all tests pass
- Commit and push changes
- Merge feat/bloop-client-bubbletea-tui → develop
- Merge develop → main
- Push all changes to origin

## Open Questions

- Should vim-style navigation (j/k) be supported in addition to arrows?
- What specific color palette should be used (brand colors or accessible standard)?
- Should API verification be mandatory or skippable?
- How should very long tunnel lists be handled (pagination, scroll, virtual list)?
- Should keyboard shortcuts for power users be added (e.g., 't' to jump to tunnels)?
- Should clipboard operations be supported for copying output?

## Complexity Tracking

No constitution violations currently require justification.
