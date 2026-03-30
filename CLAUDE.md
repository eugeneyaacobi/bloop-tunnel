# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

bloop-tunnel is a secure HTTP tunneling client that exposes local services to the internet via a WebSocket relay. It provides an interactive TUI for setup, supports YAML and environment-based configuration, and includes optional Docker service discovery.

## Development Commands

### Build and Run

```bash
# Build all packages
make build

# Run the bloop-tunnel client (from source)
make bloop-tunnel

# Run the relay server (if present in codebase)
make relay
```

### Testing

```bash
# Run all unit tests
make test
# or directly:
go test ./...

# Run specific package tests
go test ./internal/client/...
go test ./internal/client/setup/...

# Run tests with coverage
go test -cover ./...
```

### Verification

```bash
# Run full verification: fmt, tests, and Docker smoke tests
make verify

# Run Docker smoke tests only
make smoke-docker
```

### Linting

```bash
# Format code
make fmt
```

### Running the Client

```bash
# Interactive setup wizard
bloop-tunnel setup

# Run with config file
bloop-tunnel run --config ~/.bloop-tunnel/config.yaml

# Run with environment variables only
bloop-tunnel run

# Show version
bloop-tunnel run --version
```

## Architecture

### Core Components

1. **Protocol Layer** (`internal/protocol/`)
   - Defines WebSocket message types for client-relay communication
   - Message types: `ClientHello`, `RegisterTunnels`, `RequestBegin`, `ResponseBegin`, etc.
   - Uses JSON envelope format with `Type` and `Payload` fields

2. **Transport Layer** (`internal/transport/`)
   - WebSocket connection management via `gorilla/websocket`
   - Handles ping/pong keepalive (60s read deadline)
   - 10MB message size limit

3. **Client Session** (`internal/client/`)
   - `session.go`: Main client connection to relay, handles tunnel registration
   - `forward_http.go`: HTTP request forwarding from relay to local services
   - `reconnect.go`: Automatic reconnection with exponential backoff
   - `register.go`: Tunnel registration protocol
   - `enrollment.go`: Optional runtime enrollment with control plane
   - `runtime_ingest.go`: Runtime status telemetry to control plane
   - `accessmap.go`: Access control for tunnels (public, basic_auth, token_protected)

4. **Configuration** (`internal/config/`)
   - Three-layer loading: defaults → YAML file → environment variables
   - Environment variables use `BLOOP_` prefix (e.g., `BLOOP_TUNNELS_0_NAME`)
   - Indexed tunnel config via env vars: `BLOOP_TUNNELS_<n>_<FIELD>`
   - Production defaults: `api.bloop.to` for control plane, `relay.bloop.to` for relay

5. **TUI Setup** (`internal/client/tui/`, `internal/client/setup/`)
   - Built with `charmbracelet/bubbletea` and `lipgloss`
   - Interactive wizard for config generation
   - Outputs: YAML files, `.env` files, or Docker Compose blocks
   - Optional Docker container discovery (opt-in)

6. **Docker Discovery** (`internal/client/dockerdiscover/`)
   - Lists running containers with exposed ports
   - Generates candidate addresses like `container-name:3000`
   - Filter functionality for container selection

### Request Flow

1. Client connects to relay via WebSocket with `ClientHello` (auth token)
2. Client sends `RegisterTunnels` with tunnel definitions
3. Relay forwards HTTP requests as `RequestBegin` → `RequestBody` → `RequestEnd`
4. Client forwards to local service, sends back `ResponseBegin` → `ResponseBody` → `ResponseEnd`
5. Ping/pong every 20s for connection health

### Tunnel Access Modes

- `public`: No authentication required
- `basic_auth`: Username/password protection
- `token_protected`: Bearer token required

### Key Design Decisions

- **WebSocket protocol**: Chosen for bidirectional messaging and connection efficiency
- **Environment variable config**: Enables Docker-native deployment without file mounts
- **TUI-first setup**: Interactive wizard reduces config file editing complexity
- **Production defaults**: New configs default to hosted infrastructure, not localhost
- **Opt-in Docker discovery**: No automatic exposure; user confirms each tunnel

## Go Version

Requires Go 1.24.0+

## Configuration Locations

- Default config path: `~/.bloop-tunnel/config.yaml`
- Environment prefix: `BLOOP_`
- Tunnel index prefix: `BLOOP_TUNNELS_`
