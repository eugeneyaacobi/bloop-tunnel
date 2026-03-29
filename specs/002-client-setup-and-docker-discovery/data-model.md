# Data Model: bloop-client interactive setup, production defaults, and Docker endpoint discovery

## Interactive Setup Session

Represents the temporary in-memory state of a guided configuration session.

### Fields
- `config_path`: target file path for reading/writing config
- `mode`: `interactive`, `non_interactive`, or `mixed`
- `output_mode`: `yaml_file`, `env_file`, `compose_block`, or `stdout_preview`
- `control_plane_url`: chosen or default control-plane URL
- `relay_url`: chosen or default relay URL
- `auth_token_source`: `inline`, `env`, or `external_existing`
- `reconnect`: reconnect settings
- `tunnels`: ordered list of candidate tunnel definitions
- `docker_discovery_enabled`: whether discovery was requested
- `validation_errors`: collected validation failures before output

## Resolved Client Config

Represents the final runtime config after applying precedence rules.

### Fields
- `relay_url`
- `auth_token`
- `auth_token_env`
- `control_plane_url`
- `enrollment_token`
- `enrollment_token_env`
- `reconnect`
- `logging`
- `tunnels`
- `source_map`: optional metadata describing whether each field came from flags, env, file, or defaults

## Environment Tunnel Definition

Represents a single tunnel assembled from indexed environment variables.

### Fields
- `index`: numeric position from env key prefix
- `name`
- `hostname`
- `local_addr`
- `access`
- `token`
- `token_env`
- `basic_auth_username`
- `basic_auth_password`
- `basic_auth_password_env`
- `valid`: boolean after validation
- `errors`: validation failures for this entry

### Encoding
Preferred indexed format:
- `BLOOP_TUNNELS_0_NAME`
- `BLOOP_TUNNELS_0_HOSTNAME`
- `BLOOP_TUNNELS_0_LOCAL_ADDR`
- `BLOOP_TUNNELS_0_ACCESS`
- `BLOOP_TUNNELS_0_BASIC_AUTH_USERNAME`
- `BLOOP_TUNNELS_0_BASIC_AUTH_PASSWORD_ENV`

## Docker Discovery Candidate

Represents a container-backed service that could become a tunnel.

### Fields
- `container_id`
- `container_name`
- `image`
- `port`
- `protocol`
- `local_addr_candidate`
- `labels`: selected container metadata
- `published`: whether port is published externally
- `confidence`: heuristic confidence that this is an HTTP-like service worth showing
- `excluded_reason`: why the candidate was filtered or de-prioritized

## Frontend Setup Guidance

Represents UI-facing configuration guidance rather than runtime state.

### Fields
- `setup_paths`: available setup modes (`interactive_cli`, `env_only_docker`, `config_file`)
- `default_control_plane_url`
- `default_relay_url`
- `discovery_supported`: whether Docker discovery is supported/documented
- `discovery_copy`: explanatory text emphasizing opt-in confirmation
- `examples`: links or snippets shown to operators
