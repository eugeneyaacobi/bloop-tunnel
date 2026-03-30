# bloop-tunnel

Securely expose local services to the internet with zero configuration. **bloop-tunnel** creates encrypted HTTPS tunnels from your machine to `bloop.to` with automatic TLS, no port forwarding required.

## Quick Start

The **terminal UI (TUI)** is the recommended way to configure bloop-tunnel:

```bash
# Install
curl -fsSL https://install.bloop.to/install.sh | sh

# Launch the interactive terminal setup wizard
bloop-tunnel setup

# Start tunneling
bloop-tunnel run
```

The `bloop-tunnel setup` wizard walks you through everything:
- Choose production or local infrastructure
- Configure auth tokens
- Define tunnels for your services
- Optional Docker discovery (pick running containers to expose)
- Generate YAML config, `.env` file, or Docker Compose block

Your local service is now available at `https://your-app.bloop.to`.

---

## Why bloop-tunnel?

**Expose localhost securely**
- No router configuration, no port forwarding, no DNS changes
- Automatic HTTPS with valid certificates
- Works from behind NATs, firewalls, and corporate proxies

**Works where you work**
- Local development on macOS/Linux/Windows
- Docker and container deployments
- CI/CD pipelines
- Production workloads

**Simple but powerful**
- One-command setup with interactive wizard
- Configure via YAML file or environment variables
- Optional Docker service discovery
- Production-ready defaults pointing to hosted infrastructure

---

## Installation

### Native (macOS, Linux, Windows)

```bash
# Install script (recommended)
curl -fsSL https://install.bloop.to/install.sh | sh

# Manual download
wget https://github.com/eugeneyaacobi/bloop-tunnel/releases/latest/download/bloop-tunnel-linux-amd64.tar.gz
tar -xzf bloop-tunnel-linux-amd64.tar.gz
sudo mv bloop-tunnel /usr/local/bin/
```

### Docker

```bash
docker run -d \
  --name bloop-tunnel \
  -e BLOOP_CONTROL_PLANE_URL=https://api.bloop.to \
  -e BLOOP_TUNNELS_0_NAME=app \
  -e BLOOP_TUNNELS_0_HOSTNAME=app.bloop.to \
  -e BLOOP_TUNNELS_0_LOCAL_ADDR=host.docker.internal:3000 \
  ghcr.io/bloop/bloop-tunnel:latest
```

### Go Install

```bash
go install github.com/eugeneyaacobi/bloop-tunnel/cmd/bloop-tunnel@latest
```

---

## Configuration

**Terminal UI (Recommended)** — The `bloop-tunnel setup` wizard is the easiest way to configure everything. It handles:
- Control plane and relay endpoints (production defaults provided)
- Auth token configuration
- Tunnel definitions
- **Docker service discovery** — list running containers, see their ports, and add tunnels with one keypress
- Output format selection (YAML file, `.env`, or Docker Compose block)

```bash
bloop-tunnel setup
```

### Config File (Manual)

Create `~/.bloop-tunnel/config.yaml`:

```yaml
control_plane_url: https://api.bloop.to
relay_url: wss://relay.bloop.to/connect
auth_token_env: BLOOP_CLIENT_TOKEN

tunnels:
  - name: webapp
    hostname: webapp.bloop.to
    local_addr: 127.0.0.1:3000
    access: public

  - name: admin
    hostname: admin.bloop.to
    local_addr: 127.0.0.1:4000
    access: basic_auth
    basic_auth:
      username: admin
      password_env: ADMIN_PASSWORD
```

Run with config file:

```bash
bloop-tunnel run --config ~/.bloop-tunnel/config.yaml
```

### Environment Variables (Docker/Container)

```bash
export BLOOP_CONTROL_PLANE_URL=https://api.bloop.to
export BLOOP_RELAY_URL=wss://relay.bloop.to/connect
export BLOOP_AUTH_TOKEN_ENV=BLOOP_CLIENT_TOKEN

export BLOOP_TUNNELS_0_NAME=app
export BLOOP_TUNNELS_0_HOSTNAME=app.bloop.to
export BLOOP_TUNNELS_0_LOCAL_ADDR=127.0.0.1:3000
export BLOOP_TUNNELS_0_ACCESS=public

export BLOOP_TUNNELS_1_NAME=admin
export BLOOP_TUNNELS_1_HOSTNAME=admin.bloop.to
export BLOOP_TUNNELS_1_LOCAL_ADDR=127.0.0.1:4000
export BLOOP_TUNNELS_1_ACCESS=basic_auth
export BLOOP_TUNNELS_1_BASIC_AUTH_USERNAME=admin
export BLOOP_TUNNELS_1_BASIC_AUTH_PASSWORD_ENV=ADMIN_PASSWORD

bloop-tunnel run
```

### Docker Compose

```yaml
version: '3.8'
services:
  app:
    image: your-app:latest
    ports:
      - "3000:3000"

  bloop-tunnel:
    image: ghcr.io/bloop/bloop-tunnel:latest
    depends_on:
      - app
    environment:
      BLOOP_CONTROL_PLANE_URL: https://api.bloop.to
      BLOOP_RELAY_URL: wss://relay.bloop.to/connect
      BLOOP_AUTH_TOKEN_ENV: BLOOP_CLIENT_TOKEN
      BLOOP_TUNNELS_0_NAME: app
      BLOOP_TUNNELS_0_HOSTNAME: app.bloop.to
      BLOOP_TUNNELS_0_LOCAL_ADDR: app:3000
      BLOOP_TUNNELS_0_ACCESS: public
    env_file:
      - .env
```

---

## Features

### Tunnel Types

**Public**
```yaml
access: public
```
Accessible to anyone with the URL.

**Basic Auth**
```yaml
access: basic_auth
basic_auth:
  username: admin
  password_env: ADMIN_PASSWORD
```
Protect with username/password.

**Token Protected**
```yaml
access: token_protected
token_env: TUNNEL_TOKEN
```
Protect with bearer token.

### Docker Service Discovery (Terminal UI)

Discover running Docker containers directly from the `bloop-tunnel setup` wizard:

- **Lists all containers** with their images and exposed ports
- **Shows candidate addresses** like `container-name:3000`
- **One-key tunnel creation** — select a container and press Enter to add it as a tunnel
- **No auto-exposure** — you choose what gets tunneled

```bash
bloop-tunnel setup --discover-docker
```

### Production Defaults

New configurations default to:
- Control plane: `https://api.bloop.to`
- Relay: `wss://relay.bloop.to/connect`

Override for local development:
```bash
bloop-tunnel setup --control-plane-url http://localhost:8080 --relay-url ws://localhost:9000
```

---

## Commands

```bash
# Interactive setup wizard
bloop-tunnel setup [flags]

# Run tunnel client
bloop-tunnel run [flags]

# Show version
bloop-tunnel version
```

### Setup Flags

```bash
--config <path>                # Write YAML config to path (or '-' for stdout)
--output <mode>               # Output mode: yaml|env-file|compose-block
--non-interactive             # Generate output without prompts
--control-plane-url <url>     # Override control plane URL
--relay-url <url>             # Override relay URL
--discover-docker             # Offer Docker service discovery
```

### Run Flags

```bash
--config <path>               # Path to config file (or use env vars)
--control-plane-url <url>     # Override control plane URL
--relay-url <url>             # Override relay URL
--version                     # Show version and exit
```

---

## Documentation

- **[Client Installation Guide](docs/CLIENT_INSTALL.md)** — Detailed installation and configuration
- **[TUI Usage Guide](docs/TUI_USAGE.md)** — Interactive wizard keyboard shortcuts and navigation
- **[Examples](deploy/examples/)** — Sample configs and Compose files

---

## Development

### Build from Source

```bash
# Clone
git clone https://github.com/eugeneyaacobi/bloop-tunnel.git
cd bloop-tunnel

# Build
go build -o dist/bloop-tunnel ./cmd/bloop-tunnel

# Run
./dist/bloop-tunnel setup
```

### Cross-Platform Build

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dist/bloop-tunnel-darwin-arm64 ./cmd/bloop-tunnel

# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o dist/bloop-tunnel-linux-amd64 ./cmd/bloop-tunnel

# Windows (x86_64)
GOOS=windows GOARCH=amd64 go build -o dist/bloop-tunnel-windows-amd64.exe ./cmd/bloop-tunnel
```

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test ./test/integration/...

# All tests with coverage
go test -cover ./...
```

---

## Security

- **Encrypted TLS**: All traffic uses end-to-end encryption
- **Token-based Auth**: Relay connections use bearer tokens
- **No Auto-Exposure**: Docker discovery requires explicit confirmation
- **Secret Protection**: Tokens masked in TUI, never logged
- **File Permissions**: Config files written with restrictive permissions (0o600)

---

## Production Defaults

The interactive setup wizard and example configs use production defaults:

| Setting | Production Default | Local Dev Override |
|----------|-------------------|-------------------|
| Control Plane | `https://api.bloop.to` | `http://localhost:8080` |
| Relay | `wss://relay.bloop.to/connect` | `ws://localhost:9000` |
| Auth Token Env | `BLOOP_CLIENT_TOKEN` | Any env var name |

---

## Environment Variables

All configuration can be set via environment variables:

| Variable | Required | Description | Default |
|-----------|-----------|-------------|----------|
| `BLOOP_CONTROL_PLANE_URL` | No | Control plane API endpoint | `https://api.bloop.to` |
| `BLOOP_RELAY_URL` | No | Relay WebSocket endpoint | `wss://relay.bloop.to/connect` |
| `BLOOP_AUTH_TOKEN_ENV` | No | Env var containing relay auth token | `BLOOP_CLIENT_TOKEN` |
| `BLOOP_ENROLLMENT_TOKEN` | No | Enrollment token for runtime registration | — |
| `BLOOP_ENROLLMENT_TOKEN_ENV` | No | Env var containing enrollment token | — |

### Indexed Tunnel Variables

Define multiple tunnels using indexed variables:

```bash
BLOOP_TUNNELS_<n>_NAME        # Required: Tunnel name
BLOOP_TUNNELS_<n>_HOSTNAME    # Optional: Custom hostname (defaults to name.bloop.to)
BLOOP_TUNNELS_<n>_LOCAL_ADDR   # Required: Local address (host:port)
BLOOP_TUNNELS_<n>_ACCESS      # Required: public, basic_auth, or token_protected
BLOOP_TUNNELS_<n>_BASIC_AUTH_USERNAME           # If access=basic_auth
BLOOP_TUNNELS_<n>_BASIC_AUTH_PASSWORD_ENV        # If access=basic_auth
BLOOP_TUNNELS_<n>_TOKEN_ENV                    # If access=token_protected
```

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/eugeneyaacobi/bloop-tunnel/issues)
- **Discussions**: [GitHub Discussions](https://github.com/eugeneyaacobi/bloop-tunnel/discussions)

---

**bloop-tunnel** — Secure tunnels from localhost to the internet, one command away.
