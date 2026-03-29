# Tasks: BubbleTea TUI for bloop-client Setup

**Input**: Design documents from `/specs/003-bubbletea-tui-for-client-setup/`
**Prerequisites**: plan.md, spec.md

**Tests**: Unit tests for model state transitions, integration tests for full setup flows, manual testing for terminal compatibility

**Organization**: Tasks are grouped by implementation phases with clear checkpoints, enabling incremental delivery and independent testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1-US6)
- Include exact file paths in descriptions

---

## Phase 1: Foundations (Blocking Prerequisites)

**Purpose**: Set up BubbleTea infrastructure, basic model/view scaffolding, and reusable components

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Dependencies and Setup

- [ ] T001 Add BubbleTea and Lipgloss dependencies to go.mod in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/go.mod`
- [ ] T002 Create `internal/client/tui` package structure in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/`
- [ ] T003 [P] Create subdirectories for screens, models, api, and util in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/`

### Core Types and Messages

- [ ] T004 Define shared types and messages in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/types.go` (State, Config, ScreenTransitionMsg, APITestMsg, etc.)
- [ ] T005 Define keyboard event handlers in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/util/keys.go`

### Styling System

- [ ] T006 Implement Lipgloss style definitions in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/styles.go` (colors, component styles, layout helpers)

### Reusable Models

- [ ] T007 [P] Implement input field model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/input_field.go` (text and password field types)
- [ ] T008 [P] Implement select field model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/select_field.go` (dropdown/radio selection)
- [ ] T009 [P] Implement list view model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/list_view.go` (scrollable list with selection)
- [ ] T010 [P] Implement status/spinner model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/status.go` (loading indicators)

### Main Model and Entry Point

- [ ] T011 Implement TUI entry point and MainModel in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (Init, Update, View, program.Run)

### Unit Tests for Foundations

- [ ] T012 [P] Add tests for input field model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/input_field_test.go`
- [ ] T013 [P] Add tests for select field model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/select_field_test.go`
- [ ] T014 [P] Add tests for list view model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/models/list_view_test.go`

**Checkpoint**: Foundation ready - TUI infrastructure complete, all models tested, screen implementation can begin

---

## Phase 2: User Story 1 - Navigate setup through screen-based TUI (Priority: P1) 🎯 MVP

**Goal**: Implement core screens and navigation for a complete setup flow

**Independent Test**: Run TUI, navigate through all screens using keyboard, verify all configuration options accessible

### Welcome Screen

- [ ] T015 [P] [US1] Implement Welcome screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/welcome.go`
- [ ] T016 [P] [US1] Implement Welcome screen view and navigation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/welcome.go`

### Config Screen

- [ ] T017 [P] [US1] Implement Config screen model with all fields in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go`
- [ ] T018 [US1] Implement Config screen validation logic in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go` (URLs non-empty, delays positive, at least one token specified)
- [ ] T019 [P] [US1] Implement Config screen view in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go`

### Tunnels List Screen

- [ ] T020 [P] [US1] Implement Tunnels List screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnels.go`
- [ ] T021 [P] [US1] Implement Tunnels List screen view in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnels.go`

### Tunnel Form Screen

- [ ] T022 [P] [US1] Implement Tunnel Form screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnel_form.go`
- [ ] T023 [US1] Implement Tunnel Form screen validation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnel_form.go` (name/local_addr required, conditional fields based on access type)
- [ ] T024 [P] [US1] Implement Tunnel Form screen view in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnel_form.go`

### Screen Navigation

- [ ] T025 [US1] Implement screen navigation logic in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (forward/back navigation, screen stack management)
- [ ] T026 [US1] Implement quit confirmation modal in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/quit_confirm.go`

### Unit Tests for US1 Screens

- [ ] T027 [P] [US1] Add tests for Config screen state transitions in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config_test.go`
- [ ] T028 [P] [US1] Add tests for Tunnel Form validation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/tunnel_form_test.go`
- [ ] T029 [P] [US1] Add tests for screen navigation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main_test.go`

**Checkpoint**: Core screens complete - User can navigate Welcome → Config → Tunnels → Edit tunnels with keyboard

---

## Phase 3: User Story 2-4 - API Verification Screens (Priority: P1)

**Goal**: Implement verification flows for control plane connectivity, enrollment token, and relay auth token

**Independent Test**: Configure endpoints, run verification, verify success/failure feedback with clear errors

### API Verification Logic

- [ ] T030 [P] [US2] Implement control plane connectivity test in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/connectivity.go` (DNS, TLS, HTTP health check)
- [ ] T031 [P] [US3] Implement enrollment token verification in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/enrollment.go` (POST /api/runtime/enroll, parse response)
- [ ] T032 [P] [US4] Implement relay auth token verification in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/relay.go` (WebSocket connection, auth verification)

### Verification Screen

- [ ] T033 [P] [US2-US4] Implement Verification screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/verification.go`
- [ ] T034 [US2-US4] Implement Verification screen view with tabs/sections in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/verification.go`
- [ ] T035 [US2-US4] Implement Verification screen retry logic and error handling in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/verification.go`

### Integration with Config Screen

- [ ] T036 [US2] Integrate connectivity test into Config screen in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go` (inline result display)
- [ ] T037 [US2-US4] Wire Verification screen into navigation flow in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (Config → Verification → Review)

### Unit Tests for API Verification

- [ ] T038 [P] [US2] Add tests for connectivity test with mocked responses in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/connectivity_test.go`
- [ ] T039 [P] [US3] Add tests for enrollment verification in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/enrollment_test.go`
- [ ] T040 [P] [US4] Add tests for relay verification in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/api/relay_test.go`

**Checkpoint**: API verification complete - All three verification flows tested and integrated

---

## Phase 4: User Story 5 - Docker Discovery Integration (Priority: P2)

**Goal**: Integrate Docker discovery as opt-in screen with candidate selection

**Independent Test**: Mount Docker socket, discover containers, select candidates, verify added as tunnels

### Docker Discovery Screen

- [ ] T041 [P] [US5] Implement Docker Discovery screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery.go`
- [ ] T042 [P] [US5] Implement Docker Discovery screen view with candidate list in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery.go`
- [ ] T043 [US5] Implement Docker candidate selection and confirmation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery.go`
- [ ] T044 [US5] Implement graceful error handling for unavailable socket in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery.go`

### Integration with Existing Docker Discovery

- [ ] T045 [US5] Integrate existing dockerdiscover package with TUI screen in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery.go`
- [ ] T046 [US5] Add Docker Discovery toggle in Config screen in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go`
- [ ] T047 [US5] Wire Docker Discovery screen into navigation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (Tunnels → Docker Discovery → Review, optional skip)

### Unit Tests for Docker Discovery

- [ ] T048 [P] [US5] Add tests for Docker Discovery screen with mocked client in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/docker_discovery_test.go`

**Checkpoint**: Docker discovery complete - Optional screen integrated with existing discovery package

---

## Phase 5: User Story 5-6 - Review and Output Screens (Priority: P2)

**Goal**: Implement review summary and output mode selection

**Independent Test**: Complete configuration, view review summary, select output mode, verify correct output format

### Review Screen

- [ ] T049 [P] [US5] Implement Review screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/review.go`
- [ ] T050 [P] [US5] Implement Review screen view with summary sections in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/review.go`
- [ ] T051 [US5] Implement Review screen edit shortcuts in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/review.go` (jump to Config, Tunnels, Output)

### Output Screen

- [ ] T052 [P] [US6] Implement Output screen model in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output.go`
- [ ] T053 [P] [US6] Implement Output screen view with mode selection in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output.go`
- [ ] T054 [US6] Implement file output with 0o600 permissions in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output.go`
- [ ] T055 [US6] Implement stdout output for env-file and compose-block in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output.go`

### Integration with Existing Output Logic

- [ ] T056 [US5-US6] Reuse existing output writers from setup package in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output.go` (renderEnv, YAML marshaling)
- [ ] T057 [US5-US6] Wire Review and Output screens into navigation in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (Verification/Discovery → Review → Output)

### Unit Tests for Review and Output

- [ ] T058 [P] [US5] Add tests for Review screen in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/review_test.go`
- [ ] T059 [P] [US6] Add tests for Output screen file permissions and formats in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/output_test.go`

**Checkpoint**: Review and Output complete - Full setup flow from Welcome to Output working

---

## Phase 6: Integration and Polish

**Purpose**: Integrate TUI with existing setup command, handle edge cases, polish UX

### Integration with Setup Command

- [ ] T060 [US1] Modify setup command routing in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/cmd/bloop-client/main.go` to call TUI when interactive and TTY available
- [ ] T061 [US1] Implement TTY detection and fallback in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/setup/setup.go` (use existing prompts when no TTY)
- [ ] T062 [US1] Preserve all existing setup flags in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/setup/setup.go` (--non-interactive, --output, --config, etc.)

### Config Loading and Existing Config Editing

- [ ] T063 [US1] Integrate existing config loading with TUI in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/config.go` (load and preload values)
- [ ] T064 [US1] Implement "Edit Existing" flow in Welcome screen in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/welcome.go`

### UX Polish

- [ ] T065 [P] [US1] Implement terminal resize handling in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/main.go` (WindowSizeMsg)
- [ ] T066 [P] [US1] Add loading states and spinners for async operations in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/` (using status model)
- [ ] T067 [P] [US1] Add consistent footer with keyboard shortcuts in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/` (using styles)

### Testing for Integration

- [ ] T068 [P] Add integration test for full TUI setup flow in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/tests/tui/integration_test.go`
- [ ] T069 [P] Verify all existing setup tests still pass in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/setup/setup_test.go` (backward compatibility)

**Checkpoint**: TUI fully integrated - All screens working, existing tests passing, flags preserved

---

## Phase 7: Security Review

**Purpose**: Audit token handling, API verification, and file permissions; implement fixes

- [ ] T070 [Security] Spawn security review subagent to audit TUI implementation
- [ ] T071 [Security] Review auth token masking in input fields (password field type in models)
- [ ] T072 [Security] Review auth token display (ensure no plaintext in views)
- [ ] T073 [Security] Review config file permissions (verify 0o600 on write)
- [ ] T074 [Security] Review API verification timeouts and URL validation (prevent SSRF)
- [ ] T075 [Security] Review Docker discovery confirmation (ensure opt-in with explicit confirm)
- [ ] T076 [Security] Review logging and error messages (ensure no secrets leaked)
- [ ] T077 [Security] Implement any findings from security review

**Checkpoint**: Security review complete - All findings addressed

---

## Phase 8: Documentation Updates

**Purpose**: Update docs with TUI usage guide, screenshots, and troubleshooting

### Documentation Files

- [ ] T078 [P] Update README.md with TUI screenshots and quick start in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/README.md`
- [ ] T079 [P] Create TUI usage guide with keyboard shortcuts in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/docs/TUI_USAGE.md`
- [ ] T080 [P] Add API verification troubleshooting section in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/docs/TUI_USAGE.md`
- [ ] T081 [P] Update CLIENT_INSTALL.md with TUI setup instructions in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/docs/CLIENT_INSTALL.md`
- [ ] T082 [P] Update example config with production defaults in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/deploy/examples/client.example.yaml`

### Code Documentation

- [ ] T083 [P] Add godoc comments to all exported TUI functions and types in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/`
- [ ] T084 [P] Add comments for complex business logic (validation, screen transitions) in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/internal/client/tui/screens/`

**Checkpoint**: Documentation complete - All user-facing docs updated

---

## Phase 9: Final Testing and Validation

**Purpose**: Comprehensive testing across terminals, manual validation, bug fixes

### Manual Testing

- [ ] T085 [P] Test TUI on macOS Terminal and iTerm2
- [ ] T086 [P] Test TUI on Linux gnome-terminal and xterm
- [ ] T087 [P] Test TUI on Windows Terminal and PowerShell
- [ ] T088 [P] Test with different terminal sizes (80x24, 120x40)
- [ ] T089 [P] Test with light and dark terminal backgrounds
- [ ] T090 [P] Test keyboard-only navigation accessibility

### Automated Testing

- [ ] T091 [P] Run full test suite: `go test ./...` in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/`
- [ ] T092 [P] Verify all existing tests pass (backward compatibility)
- [ ] T093 [P] Verify TUI model tests pass
- [ ] T094 [P] Verify integration tests pass

### Bug Fixes

- [ ] T095 Fix any bugs discovered during testing
- [ ] T096 [P] Add regression tests for any bugs fixed

**Checkpoint**: All tests passing, manual testing complete across platforms

---

## Phase 10: Merge to Main

**Purpose**: Merge feature branch to develop and main, push to origin

### Git Operations

- [ ] T097 Run `go test ./...` to verify all tests pass
- [ ] T098 Run `go mod tidy` to clean up dependencies
- [ ] T099 Commit all changes with descriptive message
- [ ] T100 Push to origin: `git push origin feat/bloop-client-bubbletea-tui`
- [ ] T101 Merge feat/bloop-client-bubbletea-tui → develop
- [ ] T102 Push develop to origin
- [ ] T103 Merge develop → main
- [ ] T104 Push main to origin
- [ ] T105 Delete feature branch (optional)

**Checkpoint**: Feature merged to main and deployed

---

## Phase 11: Frontend Handoff

**Purpose**: Provide clear prompt for any remaining frontend work

- [ ] T106 Document frontend handoff prompt in `/root/.openclaw/workspace/worktree-bloop-tclient-bubbletea-tui/specs/003-bubbletea-tui-for-client-setup/frontend-handoff-prompt.md`

**Checkpoint**: Frontend handoff ready

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Foundations)**: No dependencies - can start immediately
- **Phase 2 (US1)**: Depends on Phase 1 completion
- **Phase 3 (US2-US4)**: Depends on Phase 1 completion; can overlap with late Phase 2
- **Phase 4 (US5 - Docker)**: Depends on Phase 1 and Phase 2 (Tunnels screen)
- **Phase 5 (US5-US6 - Review/Output)**: Depends on Phase 1, Phase 2, and Phase 3
- **Phase 6 (Integration)**: Depends on all screen phases (2-5)
- **Phase 7 (Security)**: Depends on Phase 6
- **Phase 8 (Documentation)**: Can start in parallel with Phase 6
- **Phase 9 (Testing)**: Depends on Phase 6, Phase 7, Phase 8
- **Phase 10 (Merge)**: Depends on Phase 9
- **Phase 11 (Frontend Handoff)**: Depends on Phase 10

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 1; independent implementation of core screens
- **US2 (P1)**: Depends on Phase 1 and Phase 2 (Config screen)
- **US3 (P1)**: Depends on Phase 1 and Phase 2 (Config screen)
- **US4 (P1)**: Depends on Phase 1 and Phase 2 (Config screen)
- **US5 (P2)**: Depends on Phase 1, Phase 2 (Tunnels), Phase 3 (Verification)
- **US6 (P2)**: Depends on Phase 1, Phase 5 (Review/Output)

### Parallel Opportunities

**Within Phase 1**:
- T007, T008, T009, T010 can run in parallel (reusable models)
- T012, T013, T014 can run in parallel (model tests)

**Within Phase 2**:
- T015, T016, T020, T021, T022, T024 can run in parallel (screen models/views)
- T027, T028, T029 can run in parallel (screen tests)

**Within Phase 3**:
- T030, T031, T032 can run in parallel (API verification logic)
- T038, T039, T040 can run in parallel (API verification tests)

**Within Phase 4**:
- T041, T042 can run in parallel (Docker Discovery screen)
- T048 can run in parallel with implementation

**Within Phase 5**:
- T049, T050 can run in parallel (Review screen)
- T052, T053 can run in parallel (Output screen)
- T058, T059 can run in parallel (Review/Output tests)

**Within Phase 6**:
- T065, T066, T067 can run in parallel (UX polish)
- T068, T069 can run in parallel (integration tests)

**Within Phase 8**:
- T078, T079, T080, T081, T082 can run in parallel (documentation)
- T083, T084 can run in parallel (code documentation)

**Within Phase 9**:
- T085-T090 can run in parallel (manual testing on different platforms)
- T091-T094 can run in parallel (automated testing)

## Implementation Strategy

### MVP First (US1 Only)

1. Complete Phase 1: Foundations
2. Complete Phase 2: US1 (core screens and navigation)
3. **STOP and VALIDATE**: Test US1 independently
4. Demo core TUI navigation

### Incremental Delivery

1. Complete Phase 1: Foundations → Foundation ready
2. Add Phase 2: US1 → Test → Deploy (MVP!)
3. Add Phase 3: US2-US4 → Test → Deploy
4. Add Phase 4: US5 (Docker) → Test → Deploy
5. Add Phase 5: US5-US6 (Review/Output) → Test → Deploy
6. Add Phase 6: Integration → Test → Deploy
7. Add Phase 7: Security Review → Fix findings
8. Add Phase 8: Documentation → Complete
9. Add Phase 9: Testing → All tests pass
10. Add Phase 10: Merge → Feature live

### Parallel Team Strategy

With multiple developers:

1. Team completes Phase 1 together
2. Once Phase 1 is done:
   - Developer A: Phase 2 (US1 - core screens)
   - Developer B: Phase 3 (US2-US4 - API verification)
   - Developer C: Phase 4 (US5 - Docker discovery)
3. After core screens and API verification:
   - Developer A: Phase 5 (Review/Output)
   - Developer B: Phase 6 (Integration)
   - Developer C: Phase 8 (Documentation)
4. Team completes Phase 7 (Security) together
5. Team completes Phase 9-10 (Testing and Merge) together

## Notes

- [P] tasks = different files, no dependencies within same phase
- All user stories can be independently tested after implementation
- Phase 1 is critical foundation - no other work can proceed
- Security review (Phase 7) must complete before final merge
- All existing tests must continue to pass (backward compatibility)
- Manual testing required across multiple terminals and platforms
- Document any non-obvious decisions in code comments or plan.md
