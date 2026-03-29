# Feature Specification: bloop-tunnel interactive setup, production defaults, and Docker endpoint discovery

**Feature Branch**: `[002-client-setup-and-docker-discovery]`  
**Created**: 2026-03-28  
**Status**: Draft  
**Input**: User description: "Build a new feature for the bloop-tunnel client package that allows a user to easily configure new tunnels through a terminal instead of config files, default config to use api.bloop.to, allow the Docker container to use environment variables to define unlimited endpoints, optionally use the Docker socket to offer automatic exposure, and update the readmes and frontend."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure tunnels interactively from the terminal (Priority: P1)

As an operator, I want `bloop-tunnel` to walk me through creating or updating tunnel configuration in the terminal so I can set up tunnels without hand-editing YAML like it’s a punishment.

**Why this priority**: This directly addresses the main usability problem and becomes the new operator on-ramp for native and Docker usage.

**Independent Test**: Can be fully tested by running the CLI in an interactive terminal, answering prompts for one or more tunnels, saving the result, and then launching the client successfully with the generated configuration.

**Acceptance Scenarios**:

1. **Given** a user runs the interactive setup command with no existing config, **When** they answer prompts for control plane, relay, auth token source, and at least one tunnel, **Then** the CLI writes a valid client config that the runtime can start with successfully.
2. **Given** a user has an existing config, **When** they run the interactive setup command, **Then** the CLI can preserve existing values as defaults and allow the user to add, edit, or remove tunnels before writing the updated config.
3. **Given** a user prefers automation, **When** they run setup in a non-interactive or prompt-skipping mode with flags, **Then** the CLI can still generate a valid config without requiring full manual prompt flow.

---

### User Story 2 - Run the Docker client entirely from environment variables (Priority: P1)

As an operator, I want to define any number of tunnels through Docker environment variables so I can run `bloop-tunnel` in Compose or container platforms without mounting handwritten config files.

**Why this priority**: Docker is a primary install surface and currently forces file-based configuration even for simple setups.

**Independent Test**: Can be fully tested by starting the Docker client with only environment variables, defining multiple tunnels, and verifying the client resolves and registers each configured endpoint.

**Acceptance Scenarios**:

1. **Given** a container environment with core client settings and one tunnel defined through environment variables, **When** the client starts, **Then** it loads the tunnel configuration without requiring a YAML file.
2. **Given** a container environment with multiple indexed tunnel definitions, **When** the client starts, **Then** it loads all valid tunnel entries in deterministic order and registers them normally.
3. **Given** both a config file and environment variables are provided, **When** the client starts, **Then** the documented precedence rules determine whether environment variables override file values or merge into them consistently.

---

### User Story 3 - Get sane production defaults instead of localhost nonsense (Priority: P1)

As an operator, I want generated configs and examples to default to the real production control-plane host (`api.bloop.to`) and the intended relay defaults so I do not have to reverse-engineer the public deployment shape from local examples.

**Why this priority**: Bad defaults create setup friction, bad docs, and support noise immediately.

**Independent Test**: Can be fully tested by generating a new config from the CLI or examples and verifying that production-facing default values are used unless explicitly overridden.

**Acceptance Scenarios**:

1. **Given** a user runs the interactive setup flow and accepts defaults, **When** the config is written, **Then** the control plane host defaults to `https://api.bloop.to`.
2. **Given** a user starts from documented examples, **When** they inspect example configs and Compose snippets, **Then** they see production-oriented defaults rather than `localhost` placeholders for the hosted service endpoints.
3. **Given** a local development workflow still needs overrides, **When** a user supplies explicit relay or API endpoints, **Then** the client honors those overrides without forcing production endpoints.

---

### User Story 4 - Discover Docker services and offer guided exposure (Priority: P2)

As an operator, I want the client to optionally inspect the local Docker daemon and offer container-backed services for exposure so I can set up tunnels faster when my apps already run in Docker.

**Why this priority**: This is a high-leverage operator convenience feature, but it is secondary to interactive setup and env-based configuration.

**Independent Test**: Can be fully tested by mounting the Docker socket, running containers with exposed/internal ports, invoking discovery, selecting a candidate service, and verifying that a tunnel definition is created or suggested without exposing anything automatically unless confirmed.

**Acceptance Scenarios**:

1. **Given** the Docker socket is available and readable, **When** the user runs discovery-enabled setup, **Then** the CLI lists eligible running containers and candidate ports that can be exposed.
2. **Given** the user selects one or more discovered services, **When** setup finishes, **Then** the resulting config includes the chosen services as explicit tunnel definitions.
3. **Given** the Docker socket is unavailable or permission-denied, **When** the user runs the setup flow, **Then** the CLI degrades gracefully, explains that discovery is unavailable, and continues with manual entry.
4. **Given** discovery is enabled, **When** candidate services are shown, **Then** nothing is exposed automatically until the operator explicitly confirms each tunnel definition.

---

### User Story 5 - See the new configuration model reflected in docs and frontend (Priority: P2)

As an operator, I want the docs and frontend to explain and reflect terminal setup, env-defined tunnels, and optional Docker discovery so the product does not contradict itself across surfaces.

**Why this priority**: Shipping new setup paths without documentation and UI alignment is how support tickets breed in the dark.

**Independent Test**: Can be fully tested by reviewing updated docs and frontend flows to confirm they present the interactive setup, production defaults, env-based tunnel definitions, and Docker discovery behavior consistently.

**Acceptance Scenarios**:

1. **Given** a new user reads the README or install guide, **When** they follow setup instructions, **Then** they see the interactive CLI path, Docker env path, and config-file path with clear tradeoffs.
2. **Given** a user manages installations from the frontend, **When** they view setup guidance or runtime status, **Then** the UI reflects the updated configuration paths and terminology.
3. **Given** Docker discovery is supported, **When** the frontend mentions it, **Then** it is clearly described as optional and confirmation-based rather than automatic remote exposure.

---

### Edge Cases

- What happens when interactive setup is run in a non-TTY environment?
- How are partially filled indexed environment variable groups handled when one tunnel entry is malformed or incomplete?
- What is the precedence order when the same setting is supplied via CLI flags, environment variables, and config files?
- How does the client behave when Docker socket discovery is requested but Docker is not running?
- How should the client map discovered containers that expose multiple ports or only internal ports without published bindings?
- How are secrets handled in interactive setup so they are not echoed into logs or accidentally written when the user intended env-only configuration?
- How does the system avoid accidentally exposing containers that appear eligible but are clearly infrastructure components (databases, caches, brokers)?
- What happens when env-defined tunnel indexes are sparse, duplicated, or out of order?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The client MUST provide an interactive terminal setup flow for creating or updating client configuration.
- **FR-002**: The interactive setup flow MUST support configuring at least control plane URL, relay URL, auth token source, reconnect settings, and one or more tunnels.
- **FR-003**: The interactive setup flow MUST support editing existing configuration by preloading current values as prompt defaults when a config file already exists.
- **FR-004**: The client MUST support non-interactive or partially interactive setup for scripted environments using flags and documented defaults.
- **FR-005**: The client runtime MUST support loading tunnel definitions from environment variables without requiring a config file.
- **FR-006**: The client MUST support an unbounded number of environment-defined tunnels through a deterministic multi-entry encoding scheme.
- **FR-007**: The client MUST document and implement a consistent precedence model across CLI flags, environment variables, and config files.
- **FR-008**: Newly generated client config SHOULD default `control_plane_url` to `https://api.bloop.to` unless explicitly overridden.
- **FR-009**: Newly generated client config SHOULD default relay settings to the production `bloop.to` deployment values unless explicitly overridden.
- **FR-010**: The client MUST continue supporting explicit local-development overrides for control plane and relay endpoints.
- **FR-011**: The client MUST support Docker-service discovery when the Docker socket is mounted and discovery is requested.
- **FR-012**: Docker discovery MUST be opt-in and MUST require explicit operator confirmation before adding any discovered service as a tunnel definition.
- **FR-013**: The client MUST degrade gracefully when Docker discovery is unavailable, including missing socket, permissions failure, unsupported platform, or daemon unavailability.
- **FR-014**: The client MUST avoid logging secrets entered during setup or loaded from environment variables.
- **FR-015**: The repository documentation MUST explain the interactive setup flow, env-only Docker setup, and config-file workflow.
- **FR-016**: The repository documentation MUST include examples showing multi-tunnel environment-variable configuration in Docker Compose.
- **FR-017**: The frontend MUST update its copy, guidance, and configuration expectations to reflect the new setup paths and production defaults.
- **FR-018**: The frontend MUST not imply that Docker discovery exposes services automatically without operator confirmation.
- **FR-019**: The interactive setup flow SHOULD offer to output either a file-based config, an env-file template, or a Docker Compose-ready variable block depending on operator choice.
- **FR-020**: The client SHOULD validate generated or discovered tunnel definitions before writing them to disk or attempting runtime registration.

### Key Entities *(include if feature involves data)*

- **Interactive Setup Session**: The temporary in-terminal state used to collect operator choices, defaults, validation results, and output mode during guided configuration.
- **Resolved Client Config**: The fully merged runtime configuration produced after applying flags, environment variables, and config-file values according to precedence rules.
- **Environment Tunnel Definition**: A single tunnel declaration encoded in environment variables using a repeatable indexed or structured scheme.
- **Docker Discovery Candidate**: A container + port combination identified as a possible tunnel target, along with metadata used to present and confirm it safely.
- **Frontend Setup Guidance**: The UI-facing representation of supported setup paths, defaults, and explanatory copy shown to the operator.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new operator can generate a working client configuration through the terminal setup flow in under 5 minutes without manually editing YAML.
- **SC-002**: A Docker operator can launch the client with two or more tunnels defined entirely through environment variables and complete successful registration without mounting a config file.
- **SC-003**: All generated configs and primary docs default to `https://api.bloop.to` for control-plane configuration unless the operator explicitly chooses custom values.
- **SC-004**: Docker discovery presents candidate services in under 5 seconds on a typical local developer machine with fewer than 20 running containers.
- **SC-005**: In usability testing or scripted acceptance tests, no discovered service is exposed unless the operator explicitly confirms it.
- **SC-006**: README, install guide, Docker examples, and frontend guidance all reflect the same setup model with no conflicting endpoint defaults.

## Assumptions

- `api.bloop.to` is the intended production control-plane hostname and should become the default for newly generated config and documentation.
- The hosted relay endpoint or relay bootstrap path can be represented as a stable documented default distinct from local development examples.
- The primary Docker install surface will use Compose or equivalent environments where indexed environment variables are practical.
- Docker discovery is a convenience feature for trusted local environments, not a remote orchestration feature.
- Frontend changes may live outside this repository but still need aligned specification and task tracking here.
- Existing file-based configuration remains supported for backward compatibility even after interactive and env-only paths are added.
