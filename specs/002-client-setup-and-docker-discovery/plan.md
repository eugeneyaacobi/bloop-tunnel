# Implementation Plan: bloop-tunnel interactive setup, production defaults, and Docker endpoint discovery

**Branch**: `[002-client-setup-and-docker-discovery]` | **Date**: 2026-03-28 | **Spec**: `/root/.openclaw/workspace/bloop-tunnel/specs/002-client-setup-and-docker-discovery/spec.md`
**Input**: Feature specification from `/specs/002-client-setup-and-docker-discovery/spec.md`

## Summary

Add an operator-friendly setup layer to `bloop-tunnel` that supports guided terminal configuration, env-only Docker deployments with unlimited tunnel definitions, production-ready defaults for hosted `bloop.to` infrastructure, and optional Docker socket discovery for container-backed services. Update documentation and frontend guidance so the installation flow, runtime expectations, and product copy all reflect the same configuration model.

## Technical Context

**Language/Version**: Go 1.24+  
**Primary Dependencies**: Go standard library, current config loader (`viper`), a terminal prompt library for interactive setup (or minimal stdio prompts), optional Docker client SDK for socket discovery, existing client/runtime packages  
**Storage**: File-based client config plus environment-variable resolution at runtime; no database changes required for this feature  
**Testing**: Go test, table-driven config precedence tests, CLI/setup tests, Docker discovery unit tests with mocked client responses, targeted integration tests for env-only startup  
**Target Platform**: Native macOS/Linux/Windows CLI; Docker container runtime on laptop/dev machines; frontend repo to be updated in parallel  
**Project Type**: CLI/runtime usability enhancement with cross-repo documentation/frontend impact  
**Performance Goals**: interactive setup completes in under 5 minutes; env parsing scales to at least dozens of tunnels without noticeable startup delay; Docker discovery lists candidates within 5 seconds on typical dev machines  
**Constraints**: Backward compatibility with existing YAML config; secrets must remain out of logs; Docker discovery must be opt-in and confirmation-based; frontend changes may live in a sibling repo  
**Scale/Scope**: Single-operator and small-team deployments first; unlimited env-defined tunnels by encoding rather than hardcoded count limits

## Constitution Check

### Initial Gate Review

- **Personal Infrastructure First**: Pass. This improves the self-hosted operator experience directly.
- **Small, Composable Surface Area**: Pass with caution. Interactive setup, env parsing, and Docker discovery should remain client-local concerns, not grow into a remote orchestration system.
- **Secure by Default**: Pass if discovery remains opt-in and confirmation-based, and secrets stay redacted.
- **External Edge, Simple Core**: Pass. No new edge complexity is introduced.
- **Operator-Friendly Configuration**: Strong pass. This is the main point of the feature.
- **Testable Client/Server Contracts**: Pass. Runtime config resolution and setup output can be validated independently from relay behavior.
- **Open-Source Readiness Without Premature Complexity**: Pass if Docker discovery stays narrow and well-scoped.

No constitution violations currently require justification.

## Project Structure

### Documentation (this feature)

```text
specs/002-client-setup-and-docker-discovery/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── bloop-tunnel/

internal/
├── client/
│   ├── setup/
│   └── dockerdiscover/
├── config/
├── logging/
└── auth/

deploy/
├── compose/
└── examples/

docs/
└── CLIENT_INSTALL.md
```

**Structure Decision**: Keep this feature centered in the client package. Add setup and Docker discovery support as client-local packages rather than mixing prompt logic into the runtime session code. Track frontend work in the spec/tasks even if the implementation lands in a separate repo.

## Architecture Overview

### Runtime Components

1. **Interactive Setup Command**
   - Runs as a dedicated `bloop-tunnel` subcommand (for example `setup` or `init`).
   - Collects operator input, validates it, and writes output in one of several modes.
   - Supports both new config creation and editing existing config.

2. **Resolved Config Loader**
   - Loads config from file, environment variables, and CLI flags.
   - Applies documented precedence consistently.
   - Produces one resolved in-memory config used by the existing runtime.

3. **Environment Tunnel Decoder**
   - Scans environment variables for indexed tunnel definitions.
   - Builds a deterministic ordered tunnel list.
   - Validates required per-tunnel fields before runtime registration.

4. **Docker Discovery Helper**
   - Optional client-local helper activated during setup.
   - Connects to Docker socket when available.
   - Enumerates running containers and candidate HTTP services.
   - Produces suggested tunnel definitions that still require operator confirmation.

5. **Documentation + Frontend Alignment Layer**
   - Updates docs and UI copy to match setup behavior and defaults.
   - Keeps install instructions and frontend guidance synchronized.

## CLI Design

### Proposed commands

- `bloop-tunnel run` (explicit runtime command; optional if current default main remains)
- `bloop-tunnel setup`
- `bloop-tunnel setup --config /path/client.yaml`
- `bloop-tunnel setup --print-env`
- `bloop-tunnel setup --output env-file`
- `bloop-tunnel setup --discover-docker`

### Setup flow

Prompt sequence should cover:
- control plane URL (default `https://api.bloop.to`)
- relay URL / endpoint bootstrap (production default, overridable)
- auth token source (inline for local testing vs env-based for safer setups)
- reconnect settings
- tunnel list management (add/edit/remove)
- optional Docker discovery when requested and available
- output mode:
  - write YAML config
  - print `.env` template
  - print Compose-ready env block

### Non-interactive support

The setup command should support flags that fill values without prompting, and should fail clearly when required values are missing in non-interactive mode.

## Configuration Resolution Design

### Proposed precedence

1. explicit CLI flags
2. environment variables
3. config file values
4. generated defaults

This precedence should apply to both scalar client settings and tunnel definitions where feasible.

### Environment variable scheme

Use indexed tunnel variables for unlimited tunnel count:

```text
BLOOP_CONTROL_PLANE_URL=https://api.bloop.to
BLOOP_RELAY_URL=wss://relay.bloop.to/connect
BLOOP_AUTH_TOKEN_ENV=BLOOP_CLIENT_TOKEN

BLOOP_TUNNELS_0_NAME=app
BLOOP_TUNNELS_0_HOSTNAME=app.bloop.to
BLOOP_TUNNELS_0_LOCAL_ADDR=host.docker.internal:3000
BLOOP_TUNNELS_0_ACCESS=public

BLOOP_TUNNELS_1_NAME=admin
BLOOP_TUNNELS_1_LOCAL_ADDR=host.docker.internal:4000
BLOOP_TUNNELS_1_ACCESS=basic_auth
BLOOP_TUNNELS_1_BASIC_AUTH_USERNAME=gene
BLOOP_TUNNELS_1_BASIC_AUTH_PASSWORD_ENV=BLOOP_ADMIN_PASSWORD
```

Parsing rules:
- gather all indexed prefixes matching `BLOOP_TUNNELS_<n>_`
- sort numerically by `<n>`
- ignore empty indexes only when the whole tunnel group is absent
- fail validation for partially defined tunnel groups missing required fields

### Default values

- `control_plane_url`: `https://api.bloop.to`
- relay default: production hosted endpoint documented by product decision
- reconnect defaults remain current unless product wants different operator defaults

## Docker Discovery Design

### Candidate selection

Discovery should identify likely HTTP services based on:
- running container name
- exposed/container ports
- labels when available
- exclusion heuristics for obvious infrastructure containers (db, redis, postgres, rabbitmq, etc.)

### Safety constraints

- discovery only runs when explicitly requested
- no automatic exposure
- each candidate must be confirmed or skipped
- selected candidates become ordinary explicit tunnel definitions in config output
- discovery failures must be non-fatal to the broader setup flow

## Frontend Considerations

Frontend should reflect:
- terminal setup as the primary guided setup path
- env-only Docker setup as a first-class deployment option
- production defaults pointing at `api.bloop.to`
- optional Docker discovery as a local convenience feature, not a magic auto-publish mechanism
- any installation status copy that currently assumes config-file-only setup

Because frontend code likely lives outside this repo, plan/tasks should mark those changes as cross-repo deliverables requiring Prism/front-end ownership.

## Testing Strategy

### Unit Tests
- config precedence resolution
- env tunnel decoding
- setup input validation
- Docker candidate filtering/exclusion heuristics
- redaction behavior for secrets

### CLI/Behavior Tests
- setup generates valid YAML output
- setup generates env-file/Compose output
- setup edits existing config correctly
- non-interactive mode fails cleanly when required fields are missing

### Integration Tests
- env-only startup registers multiple tunnels successfully
- mixed config file + env override behavior matches documented precedence
- Docker discovery mock path produces confirmed tunnel definitions without runtime registration side effects

### Manual Verification
- native CLI setup on macOS/Linux
- Docker Compose startup with env-only config
- Docker socket discovery with a few running app containers
- docs walkthrough matches actual commands
- frontend copy and install guidance match backend behavior

## Implementation Phases

### Phase 0 - Artifact completion
- Write research decisions for prompt library, env encoding, and Docker discovery scope.
- Write data model for setup state and resolved config.
- Write quickstart for interactive and env-only Docker setup.

### Phase 1 - Config foundation
- Refactor config loading into a resolved config pipeline.
- Add env tunnel parsing and precedence logic.
- Add tests for merge behavior and validation.

### Phase 2 - Interactive CLI setup
- Add `bloop-tunnel setup` command.
- Implement prompt flow, validation, output writers, and existing-config editing.
- Add tests for interactive and non-interactive modes.

### Phase 3 - Docker discovery
- Add optional Docker client integration and candidate selection.
- Wire discovery into setup flow with confirmation gates.
- Add mock-backed tests and graceful failure behavior.

### Phase 4 - Docs and frontend alignment
- Update README and install docs.
- Update Docker Compose examples.
- Coordinate frontend copy and setup guidance updates.

## Open Questions

- What exact hosted relay URL should become the production default alongside `https://api.bloop.to`?
- Should env tunnel definitions merge with file tunnels or replace them entirely when any env tunnel is present?
- Should `bloop-tunnel` keep backward-compatible single-command runtime behavior, or move to explicit subcommands (`setup`, `run`)?
- Which prompt library offers the best balance of portability and minimal dependency weight?
- How aggressive should Docker discovery exclusion heuristics be for common infrastructure containers?

## Complexity Tracking

No constitution violations currently require justification.
