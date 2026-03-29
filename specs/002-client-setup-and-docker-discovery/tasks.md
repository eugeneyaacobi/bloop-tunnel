# Tasks: bloop-tunnel interactive setup, production defaults, and Docker endpoint discovery

**Input**: Design documents from `/specs/002-client-setup-and-docker-discovery/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Tests**: Config resolution, setup behavior, and env-only Docker startup tests are required because operator UX and configuration correctness are the core product behavior here.

**Organization**: Tasks are grouped by foundations, user stories, and cross-repo deliverables so the client/runtime, docs, and frontend can move without drifting apart.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (`US1`, `US2`, `US3`, `US4`, `US5`)
- Include exact file paths in descriptions

---

## Phase 1: Foundations (Blocking Prerequisites)

**Purpose**: Establish the config resolution model, defaults, and output structures used by all later setup flows.

**⚠️ CRITICAL**: User-story work should not start until config precedence and runtime loading rules are settled.

- [ ] T001 Refactor client config loading into a resolved config pipeline in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client.go`
- [ ] T002 [P] Define config precedence rules and defaults in `/root/.openclaw/workspace/bloop-tunnel/specs/002-client-setup-and-docker-discovery/research.md`
- [ ] T003 [P] Add data model for setup state, resolved config, and env tunnel definitions in `/root/.openclaw/workspace/bloop-tunnel/specs/002-client-setup-and-docker-discovery/data-model.md`
- [ ] T004 [P] Add quickstart covering interactive setup, env-only Docker setup, and config-file fallback in `/root/.openclaw/workspace/bloop-tunnel/specs/002-client-setup-and-docker-discovery/quickstart.md`
- [ ] T005 [P] Implement production defaults for control plane and relay settings in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client.go`, `/root/.openclaw/workspace/bloop-tunnel/deploy/examples/client.example.yaml`, and `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md`
- [ ] T006 [P] Add unit tests for defaults and precedence behavior in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client_test.go`

**Checkpoint**: Runtime config has a stable precedence model, default endpoints, and spec artifacts documenting them.

---

## Phase 2: User Story 1 - Configure tunnels interactively from the terminal (Priority: P1) 🎯 MVP

**Goal**: Let operators create or edit tunnel configuration through a guided CLI instead of hand-editing YAML.

**Independent Test**: Run the setup command against empty and existing config files, generate valid output, and start `bloop-tunnel` successfully with the result.

### Tests for User Story 1

- [ ] T007 [P] [US1] Add setup command behavior tests for new config generation in `/root/.openclaw/workspace/bloop-tunnel/cmd/bloop-tunnel/setup_test.go`
- [ ] T008 [P] [US1] Add setup tests for editing existing config defaults in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/setup_test.go`
- [ ] T009 [P] [US1] Add non-interactive validation tests in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/validation_test.go`

### Implementation for User Story 1

- [ ] T010 [P] [US1] Add setup command/subcommand wiring in `/root/.openclaw/workspace/bloop-tunnel/cmd/bloop-tunnel/main.go`
- [ ] T011 [P] [US1] Implement setup session model and prompt flow in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/session.go`
- [ ] T012 [US1] Implement tunnel add/edit/remove prompt flow in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/tunnels.go`
- [ ] T013 [US1] Implement config output writers for YAML, env-file, and Compose-friendly output in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/output.go`
- [ ] T014 [US1] Implement existing-config preload/edit support in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/load_existing.go`
- [ ] T015 [US1] Ensure secret inputs and setup logs are redaction-safe in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/session.go` and `/root/.openclaw/workspace/bloop-tunnel/internal/logging/logging.go`

**Checkpoint**: Operators can generate or update working client config from the terminal without manual YAML editing.

---

## Phase 3: User Story 2 - Run the Docker client entirely from environment variables (Priority: P1)

**Goal**: Support unlimited environment-defined tunnel entries for Docker and Compose deployments.

**Independent Test**: Start `bloop-tunnel` with only environment variables defining multiple tunnels and verify successful config resolution and registration.

### Tests for User Story 2

- [ ] T016 [P] [US2] Add env tunnel decoding tests for indexed tunnel groups in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client_env_test.go`
- [ ] T017 [P] [US2] Add env/file precedence tests in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client_merge_test.go`
- [ ] T018 [P] [US2] Add integration test for env-only multi-tunnel startup in `/root/.openclaw/workspace/bloop-tunnel/test/integration/env_only_startup_test.go`

### Implementation for User Story 2

- [ ] T019 [P] [US2] Implement environment tunnel scanning and indexed decoding in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client_env.go`
- [ ] T020 [US2] Implement merge/override logic for env + file config in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client.go`
- [ ] T021 [US2] Update runtime startup to allow env-only config execution in `/root/.openclaw/workspace/bloop-tunnel/cmd/bloop-tunnel/main.go`
- [ ] T022 [US2] Add Docker Compose and `.env` examples for multi-tunnel env config in `/root/.openclaw/workspace/bloop-tunnel/deploy/compose/example-client-compose.yml`, `/root/.openclaw/workspace/bloop-tunnel/deploy/examples/client.env.example`, and `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md`

**Checkpoint**: Docker operators can launch bloop-tunnel with an arbitrary number of tunnels using only environment variables.

---

## Phase 4: User Story 3 - Get sane production defaults (Priority: P1)

**Goal**: Make generated configs, examples, and docs point at the real hosted defaults instead of local placeholders.

**Independent Test**: Generate new configs and inspect docs/examples to confirm production defaults are applied unless explicitly overridden.

### Tests for User Story 3

- [ ] T023 [P] [US3] Add tests asserting `https://api.bloop.to` default behavior in `/root/.openclaw/workspace/bloop-tunnel/internal/config/client_test.go`
- [ ] T024 [P] [US3] Add setup output tests confirming production default values in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/output_test.go`

### Implementation for User Story 3

- [ ] T025 [P] [US3] Update client example config to production defaults in `/root/.openclaw/workspace/bloop-tunnel/deploy/examples/client.example.yaml`
- [ ] T026 [P] [US3] Update README install and positioning copy in `/root/.openclaw/workspace/bloop-tunnel/README.md`
- [ ] T027 [US3] Update install docs and quickstart examples in `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md` and `/root/.openclaw/workspace/bloop-tunnel/specs/002-client-setup-and-docker-discovery/quickstart.md`

**Checkpoint**: Primary setup surfaces stop telling users to aim at localhost unless they explicitly choose local development.

---

## Phase 5: User Story 4 - Discover Docker services and offer guided exposure (Priority: P2)

**Goal**: Offer opt-in Docker socket discovery during setup and convert chosen services into explicit tunnel definitions.

**Independent Test**: With Docker socket mounted and mock or real containers running, discovery lists candidates, requires confirmation, and writes selected definitions without exposing anything automatically.

### Tests for User Story 4

- [ ] T028 [P] [US4] Add Docker discovery candidate filtering tests in `/root/.openclaw/workspace/bloop-tunnel/internal/client/dockerdiscover/discover_test.go`
- [ ] T029 [P] [US4] Add graceful failure tests for missing socket and permission errors in `/root/.openclaw/workspace/bloop-tunnel/internal/client/dockerdiscover/errors_test.go`
- [ ] T030 [P] [US4] Add setup-flow tests for discovery confirmation behavior in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/docker_discovery_test.go`

### Implementation for User Story 4

- [ ] T031 [P] [US4] Implement Docker discovery client and candidate model in `/root/.openclaw/workspace/bloop-tunnel/internal/client/dockerdiscover/discover.go`
- [ ] T032 [US4] Implement candidate filtering/exclusion heuristics in `/root/.openclaw/workspace/bloop-tunnel/internal/client/dockerdiscover/filter.go`
- [ ] T033 [US4] Wire discovery prompts and explicit confirmation into setup flow in `/root/.openclaw/workspace/bloop-tunnel/internal/client/setup/docker_discovery.go`
- [ ] T034 [US4] Update Docker examples and docs to explain optional socket mounting in `/root/.openclaw/workspace/bloop-tunnel/deploy/compose/example-client-compose.yml` and `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md`

**Checkpoint**: Docker-backed services can be suggested and selected safely during setup without any automatic exposure.

---

## Phase 6: User Story 5 - Reflect the new model in docs and frontend (Priority: P2)

**Goal**: Align docs and frontend guidance with the interactive setup flow, env-only Docker support, and optional Docker discovery.

**Independent Test**: Review docs and frontend changes to confirm consistent terminology, defaults, and setup paths.

### Tests / Validation for User Story 5

- [ ] T035 [P] [US5] Review README, install docs, and example files for conflicting defaults and terminology in `/root/.openclaw/workspace/bloop-tunnel/README.md`, `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md`, and `/root/.openclaw/workspace/bloop-tunnel/deploy/examples/`
- [ ] T036 [P] [US5] Review frontend setup copy and flows in the frontend repo for alignment with this spec

### Implementation for User Story 5

- [ ] T037 [P] [US5] Update README with setup-path comparison and examples in `/root/.openclaw/workspace/bloop-tunnel/README.md`
- [ ] T038 [P] [US5] Update client install guide with interactive setup, env-only Docker, and discovery caveats in `/root/.openclaw/workspace/bloop-tunnel/docs/CLIENT_INSTALL.md`
- [ ] T039 [P] [US5] Update example assets in `/root/.openclaw/workspace/bloop-tunnel/deploy/examples/` and `/root/.openclaw/workspace/bloop-tunnel/deploy/compose/`
- [ ] T040 [US5] Implement frontend copy and workflow updates in the frontend repo to reflect terminal setup, env-only configuration, production defaults, and opt-in Docker discovery

**Checkpoint**: Product surfaces tell the same story instead of each improvising jazz in a different key.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Foundations)**: starts immediately
- **Phase 2 (US1)**: depends on Phase 1
- **Phase 3 (US2)**: depends on Phase 1; can overlap with late US1 work where files do not conflict
- **Phase 4 (US3)**: depends on Phase 1 and should align with US1 output behavior
- **Phase 5 (US4)**: depends on Phase 2 foundations for setup flow and Phase 1 config model
- **Phase 6 (US5)**: depends on the chosen setup behavior being stable enough to document

### User Story Dependencies

- **US1 (P1)**: depends on config precedence and defaults from Foundations
- **US2 (P1)**: depends on config precedence and runtime loading rules from Foundations
- **US3 (P1)**: depends on Foundations and should be coordinated with US1/US2 outputs
- **US4 (P2)**: depends on US1 setup flow scaffolding
- **US5 (P2)**: depends on US1–US4 decisions being stable enough to document and reflect in UI

### Parallel Opportunities

- Foundation documentation tasks can run in parallel with config test scaffolding
- US1 setup flow internals can split across command wiring, prompt flow, and output writing
- US2 env parsing and docs/examples can run in parallel
- US4 Docker discovery client and filtering heuristics can run in parallel before setup integration
- US5 docs and frontend work can proceed in parallel once the final config model is stable

## Implementation Strategy

### MVP First (Recommended)

1. Complete Foundations
2. Ship interactive setup (US1)
3. Ship env-only Docker support (US2)
4. Fix production defaults everywhere (US3)
5. Then add Docker discovery (US4)
6. Finish docs/frontend alignment (US5)

### Incremental Delivery

1. Stable config precedence + defaults
2. Guided CLI setup
3. Unlimited env-defined tunnel support
4. Docs/examples updated to match
5. Optional Docker discovery
6. Frontend alignment

## Notes

- Keep Docker discovery explicitly optional and local-only in spirit
- Do not let env support become an undocumented side-channel; it must be first-class and tested
- Prefer additive backward-compatible changes before any CLI breaking changes
- If frontend lives in a separate repo, keep exact deliverables explicit so this spec still governs the whole feature
