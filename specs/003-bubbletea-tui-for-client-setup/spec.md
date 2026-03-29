# Feature Specification: BubbleTea TUI for bloop-tunnel Setup

**Feature Branch**: `[feat/bloop-client-bubbletea-tui]`  
**Created**: 2026-03-29  
**Status**: Draft  
**Input**: Upgrade existing bufio-based setup to BubbleTea TUI with enhanced UX, API verification, and secure token handling

## Background

The current `bloop-tunnel setup` command uses basic bufio prompts that provide minimal UX. This spec upgrades the interactive setup to use BubbleTea (Charmbracelet's Elm-architecture TUI framework) for a modern, responsive terminal interface with:

- Screen-based navigation
- Form validation with immediate feedback
- API verification flows
- Better accessibility and keyboard navigation
- Visual consistency with production-grade tools

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Navigate setup through screen-based TUI (Priority: P1)

As an operator, I want to walk through a polished terminal UI with clear screens and intuitive keyboard navigation so I can configure bloop-tunnel efficiently without memorizing prompt sequences.

**Why this priority**: Replaces the current linear prompt flow with a modern, discoverable interface that reduces operator error and improves first-time setup experience.

**Independent Test**: Can be fully tested by running `bloop-tunnel setup`, navigating through all screens using keyboard, and verifying that all configuration options are accessible and properly validated.

**Acceptance Scenarios**:

1. **Given** a user runs `bloop-tunnel setup` in a terminal, **When** they navigate through screens using arrow keys, tab, and enter, **Then** each screen renders correctly and transitions appropriately.
2. **Given** a user is on any screen, **When** they press `q` or `ctrl+c`, **Then** the TUI exits cleanly with a confirmation prompt if there are unsaved changes.
3. **Given** a user has made edits, **When** they press `esc` on a non-modal screen, **Then** they return to the previous screen or main menu without losing progress.

### User Story 2 - Verify control plane connectivity before enrollment (Priority: P1)

As an operator, I want to test connectivity to the control plane API before attempting enrollment so I can diagnose network or DNS issues early instead of failing cryptically during enrollment.

**Why this priority**: Prevents enrollment failures that could be avoided by validating connectivity upfront, improving operator trust in the setup process.

**Independent Test**: Can be tested by configuring an invalid control plane URL and verifying that the connectivity test fails with a clear, actionable error message before enrollment is attempted.

**Acceptance Scenarios**:

1. **Given** a user enters a control plane URL, **When** the TUI initiates a connectivity test, **Then** it shows a loading state and reports success or failure with specific error details (DNS resolution, TLS handshake, HTTP response).
2. **Given** connectivity fails, **When** the user is presented with options, **Then** they can retry with a different URL, skip verification and continue, or exit to fix their environment.
3. **Given** connectivity succeeds, **When** the TUI proceeds to the next screen, **Then** the success state is briefly displayed before advancing automatically.

### User Story 3 - Verify enrollment token before completing setup (Priority: P1)

As an operator, I want to test my enrollment token against the control plane before finishing setup so I can catch invalid or expired tokens immediately.

**Why this priority**: Invalid tokens are a common source of setup failures; validating them early prevents operators from saving a config that will fail on first run.

**Independent Test**: Can be tested by providing an invalid enrollment token and verifying that the verification step fails with a clear error message explaining the issue (invalid, expired, unauthorized).

**Acceptance Scenarios**:

1. **Given** a user has entered an enrollment token (or enrollment token env variable), **When** the TUI performs token verification, **Then** it sends a test enrollment request and displays the result.
2. **Given** token verification fails, **When** the user sees the error, **Then** the error message explains whether the token is invalid, expired, or unauthorized, and offers to retry with a new token.
3. **Given** token verification succeeds, **When** the TUI proceeds, **Then** it displays the returned installation ID and confirms that the runtime will be registered on first run.

### User Story 4 - Verify relay auth token before completing setup (Priority: P1)

As an operator, I want to test my relay auth token against the relay WebSocket endpoint so I can confirm that the client will successfully connect when starting.

**Why this priority**: Relay connection failures are difficult to debug after the fact; validating the token upfront prevents silent failures.

**Independent Test**: Can be tested by providing an invalid relay auth token and verifying that the verification step fails with a clear error message before the setup flow completes.

**Acceptance Scenarios**:

1. **Given** a user has entered a relay auth token (or auth token env variable), **When** the TUI performs relay token verification, **Then** it attempts a WebSocket connection to the relay and reports the result.
2. **Given** relay token verification fails, **When** the user sees the error, **Then** the error message indicates whether the issue is with the token, the relay URL, or network connectivity.
3. **Given** relay token verification succeeds, **When** the TUI proceeds, **Then** it confirms that the client will be able to connect to the relay on startup.

### User Story 5 - Add and manage tunnels (Priority: P1)

As an operator, I want a dedicated screen to add new tunnels by specifying service name, local IP address, and port number, and view/modify/delete existing tunnels in a management screen.

**Why this priority**: Separating tunnel configuration into dedicated screens makes it clear and easy to manage multiple tunnels, with distinct IP and Port fields for better UX.

**Independent Test**: Can be tested by navigating to the endpoints screen, adding a tunnel with name/IP/port, and verifying it appears in the tunnel management list with edit/delete options available.

**Acceptance Scenarios**:

1. **Given** a user is on the endpoints screen, **When** they enter a service name, local IP, and port number, **Then** the tunnel is added to the configuration and appears in the tunnel management list.
2. **Given** a user is on the tunnel management screen, **When** they view the list of configured tunnels, **Then** they see each tunnel with its name, local IP:port, hostname (if any), and access mode.
3. **Given** a user selects an existing tunnel from the management screen, **When** they choose to edit, **Then** they can modify the name, IP, port, and other settings, with changes saved to the configuration.
4. **Given** a user selects an existing tunnel from the management screen, **When** they choose to delete, **Then** the tunnel is removed from the configuration after confirmation.
5. **Given** a user has added or modified tunnels, **When** they proceed to the review screen, **Then** they see a complete summary of all tunnels including name, local IP:port, hostname (if any), and access mode.

### User Story 6 - Choose output mode interactively (Priority: P2)

As an operator, I want to select between YAML file, env-file, or Compose block output at the end of setup so I can get the format that matches my deployment strategy.

**Why this priority**: Different operators have different deployment preferences; making this interactive improves UX for all workflows.

**Independent Test**: Can be tested by navigating to the output selection screen and verifying that each output mode generates the correct format.

**Acceptance Scenarios**:

1. **Given** a user has completed configuration, **When** they reach the output selection screen, **Then** they can choose between YAML file, env-file, or Compose block using arrow keys and enter.
2. **Given** a user selects YAML file output, **When** they confirm, **Then** the TUI prompts for a file path and writes the config with appropriate permissions (0o600).
3. **Given** a user selects env-file or Compose block, **When** they confirm, **Then** the TUI displays the output directly in the terminal for copying or piping to a file.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The TUI MUST use BubbleTea framework (not Huh?, not Gum) following Elm-style model-view-update architecture.
- **FR-002**: The TUI MUST support screen-based navigation with the following screens: Welcome → Config entry → Endpoints (add tunnel) → Tunnels (list/manage) → Docker discovery (opt-in) → Verification → Review → Output mode selection.
- **FR-003**: The TUI MUST provide consistent keyboard navigation: arrows for selection, tab for next field/section, enter for confirm/action, esc for back/cancel, q/ctrl+c for quit.
- **FR-004**: The TUI MUST use Lipgloss for styling with a consistent color scheme that works in both light and dark terminals.
- **FR-005**: The TUI MUST display loading states for async operations (API verification, Docker discovery).
- **FR-006**: The TUI MUST provide immediate feedback for form validation errors inline with the relevant field.
- **FR-007**: The TUI MUST include an API verification screen for testing control plane connectivity.
- **FR-008**: The TUI MUST include an enrollment token verification screen that tests the token against the control plane.
- **FR-009**: The TUI MUST include a relay auth token verification screen that tests WebSocket connectivity.
- **FR-010**: The TUI MUST include a tunnel registration preview screen showing all configured tunnels before confirmation.
- **FR-011**: The TUI MUST preserve all existing setup flags (--non-interactive, --output, --config, etc.) for backward compatibility.
- **FR-012**: The TUI MUST fall back to non-interactive mode when TTY is unavailable.
- **FR-013**: The TUI MUST prompt for confirmation before exiting if unsaved changes exist.
- **FR-014**: The TUI MUST support Docker discovery as an opt-in screen that can be skipped or enabled at any time.
- **FR-015**: The TUI MUST gracefully handle all errors with clear, actionable messages and recovery options (retry, skip, exit).

### Non-Functional Requirements

- **NFR-001**: The TUI MUST render smoothly at 60fps on typical terminal sizes (80x24 minimum, 120x40 recommended).
- **NFR-002**: The TUI MUST support terminal resize events without losing state.
- **NFR-003**: The TUI MUST be accessible to screen readers by using clear text labels and avoiding ASCII-only visual cues.
- **NFR-004**: The TUI MUST not block the main thread during async operations (use TickMsg for status updates).
- **NFR-005**: The TUI MUST maintain all existing test coverage for setup behavior.
- **NFR-006**: The TUI MUST be testable with table-driven tests for model state transitions.

### Security Requirements

- **SR-001**: Auth tokens MUST never be displayed in plaintext in the TUI.
- **SR-002**: Auth tokens MUST be masked during input (password-style field).
- **SR-003**: Configuration files MUST be written with restrictive permissions (0o600).
- **SR-004**: API verification requests MUST use appropriate timeouts (5s for connectivity, 10s for enrollment/relay).
- **SR-005**: The TUI MUST redact all sensitive values from logs and error messages.
- **SR-006**: Docker discovery MUST require explicit confirmation before adding any discovered service as a tunnel.

## Key Entities

### BubbleTea Model
- **MainModel**: Root model managing screen state and navigation stack
- **WelcomeModel**: Welcome screen with quick start options
- **ConfigModel**: Form for entering control plane, relay, auth tokens, and reconnect settings
- **TunnelListModel**: List view for managing tunnels (add/edit/remove)
- **TunnelFormModel**: Form for editing a single tunnel
- **DockerDiscoveryModel**: Screen for discovering Docker services
- **VerificationModel**: Screen for running API verification tests
- **ReviewModel**: Summary screen showing all configuration
- **OutputModel**: Screen for selecting output mode and generating output
- **QuitConfirmModel**: Confirmation modal for exiting with unsaved changes

### Messages
- **TickMsg**: For status updates during async operations
- **ScreenTransitionMsg**: For navigating between screens
- **APITestMsg**: For API verification results
- **DockerDiscoveryMsg**: For Docker discovery results
- **ConfigChangedMsg**: For tracking unsaved changes
- **WindowSizeMsg**: For handling terminal resize

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new operator can complete full setup in under 5 minutes using only keyboard navigation.
- **SC-002**: All API verification flows (connectivity, enrollment, relay) complete in under 10 seconds with clear success/failure feedback.
- **SC-003**: Existing tests for setup behavior continue to pass without modification (backward compatibility).
- **SC-004**: The TUI renders correctly on macOS, Linux, and Windows terminals at 80x24 and larger.
- **SC-005**: Security review confirms no plaintext secrets are logged or displayed in the TUI.
- **SC-006**: Documentation includes TUI usage guide with screenshots, keyboard shortcuts reference, and troubleshooting guide.

## Assumptions

- BubbleTea and Lipgloss are approved dependencies in the project
- Operators have basic terminal familiarity (arrows, tab, enter, esc, q)
- The existing config model and output logic will be reused
- API verification endpoints match the documented enrollment and relay protocols
- The existing Docker discovery implementation will be integrated into the TUI

## Open Questions

- Should the TUI support vim-style navigation keys (j/k) in addition to arrows?
- What color scheme should be used for the TUI (brand colors, accessible palette, terminal-native)?
- Should API verification be mandatory or skippable with a warning?
- How should the TUI handle very long tunnel lists that don't fit on screen?
- Should the TUI support configuration templates or presets?
- How should the TUI handle keyboard shortcuts for power users (e.g., 't' to jump to tunnels)?
