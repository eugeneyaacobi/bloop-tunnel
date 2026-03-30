# bloop-tunnel Terminal UI (TUI) Guide

The `bloop-tunnel setup` command launches an interactive terminal wizard for configuring tunnels without touching a single config file.

## Launch the TUI

```bash
bloop-tunnel setup
```

Add optional flags:

```bash
--discover-docker          # Offer Docker service discovery
--control-plane-url <url>  # Override control plane URL (for local dev)
--relay-url <url>          # Override relay URL (for local dev)
```

---

## Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `Space` | Select / Confirm |
| `q` | Quit (prompts if there are unsaved changes) |
| `Ctrl+C` | Quit (prompts if there are unsaved changes) |

---

## Screens

### 1. Welcome

Choose how you want to configure bloop-tunnel:

- **New Setup** — Start fresh with production defaults
- **Edit Existing** — Load and modify an existing config file
- **Quick Start** — Skip to output with default configuration

### 2. Configuration

Configure the control plane and relay endpoints:

- **Control Plane URL** — API endpoint for tunnel registration (default: `https://api.bloop.to`)
- **Relay URL** — WebSocket endpoint for tunnel connections (default: `wss://relay.bloop.to/connect`)
- **Auth Token Env** — Environment variable name containing your auth token (default: `BLOOP_CLIENT_TOKEN`)

For local development, override the production defaults:

```bash
bloop-tunnel setup \
  --control-plane-url http://localhost:8080 \
  --relay-url ws://localhost:9000
```

### 3. Docker Discovery (Optional)

If you launch with `--discover-docker`, the TUI lists all running Docker containers:

```
[1] webapp (nginx:latest) -> webapp:80
[2] api (node:20) -> api:3000
[3] db (postgres:16) -> db:5432
```

For each container, you see:

- **Container name** — The Docker container name
- **Image** — The container image
- **Local address** — The candidate address for tunneling (e.g., `webapp:80`)

**Add a tunnel:**
1. Navigate to a container
2. Press `Enter` to select it
3. The TUI creates a tunnel config with the container's suggested name and port

**Skip discovery:**
- If no containers are found or you prefer manual entry, continue to the next screen

### 4. Tunnels

Manage your tunnel definitions:

- View all configured tunnels
- Add new tunnels (press `Enter` on "Add Tunnel")
- Edit existing tunnels (navigate and press `Enter`)
- Delete tunnels (navigate and confirm deletion)

For each tunnel, you configure:

- **Name** — Tunnel identifier (used in default hostname)
- **Hostname** — Optional custom hostname (defaults to `{name}.bloop.to`)
- **Local Address** — Local service address (`host:port`)
- **Access Control** — `public`, `basic_auth`, or `token_protected`

### 5. Review

Review your complete configuration before generating output:

- Control plane and relay settings
- Auth token configuration
- All tunnel definitions
- Any discovered Docker services

Make changes or proceed to generate config files.

### 6. Output

Choose how to save your configuration:

- **YAML File** — Write to `~/.bloop-tunnel/config.yaml` (or custom path)
- **Env File** — Generate `.env` with `BLOOP_*` environment variables
- **Docker Compose** — Generate a Compose service block

Press `Enter` to write the file and exit the TUI.

---

## Output Examples

### YAML Config

Generated at `~/.bloop-tunnel/config.yaml`:

```yaml
control_plane_url: https://api.bloop.to
relay_url: wss://relay.bloop.to/connect
auth_token_env: BLOOP_CLIENT_TOKEN

tunnels:
  - name: webapp
    hostname: webapp.bloop.to
    local_addr: webapp:80
    access: public

  - name: api
    hostname: api.bloop.to
    local_addr: api:3000
    access: basic_auth
    basic_auth:
      username: admin
      password_env: API_PASSWORD
```

Run with the config:

```bash
bloop-tunnel run
```

### Env File

Generated `.env`:

```bash
export BLOOP_CONTROL_PLANE_URL=https://api.bloop.to
export BLOOP_RELAY_URL=wss://relay.bloop.to/connect
export BLOOP_AUTH_TOKEN_ENV=BLOOP_CLIENT_TOKEN

export BLOOP_TUNNELS_0_NAME=webapp
export BLOOP_TUNNELS_0_HOSTNAME=webapp.bloop.to
export BLOOP_TUNNELS_0_LOCAL_ADDR=webapp:80
export BLOOP_TUNNELS_0_ACCESS=public

export BLOOP_TUNNELS_1_NAME=api
export BLOOP_TUNNELS_1_HOSTNAME=api.bloop.to
export BLOOP_TUNNELS_1_LOCAL_ADDR=api:3000
export BLOOP_TUNNELS_1_ACCESS=basic_auth
export BLOOP_TUNNELS_1_BASIC_AUTH_USERNAME=admin
export BLOOP_TUNNELS_1_BASIC_AUTH_PASSWORD_ENV=API_PASSWORD
```

Run with env vars:

```bash
source .env
bloop-tunnel run
```

### Docker Compose

Generated Compose block:

```yaml
bloop-tunnel:
  image: ghcr.io/bloop/bloop-tunnel:latest
  environment:
    BLOOP_CONTROL_PLANE_URL: https://api.bloop.to
    BLOOP_RELAY_URL: wss://relay.bloop.to/connect
    BLOOP_AUTH_TOKEN_ENV: BLOOP_CLIENT_TOKEN
    BLOOP_TUNNELS_0_NAME: webapp
    BLOOP_TUNNELS_0_HOSTNAME: webapp.bloop.to
    BLOOP_TUNNELS_0_LOCAL_ADDR: webapp:80
    BLOOP_TUNNELS_0_ACCESS: public
```

---

## Tips

- **Save often** — The TUI warns if you try to quit with unsaved changes
- **Test locally first** — Use `--control-plane-url` and `--relay-url` flags for local development
- **Docker discovery** — Works on Linux and macOS; requires Docker daemon running
- **Secret protection** — Auth tokens are masked in the TUI and never logged

---

## Exit and Resume

Quit at any time with `q` or `Ctrl+C`. If you've made changes, the TUI prompts you to confirm before exiting.

After completing setup, run:

```bash
bloop-tunnel run
```

To start your tunnels.

---

For more information, see [README.md](../README.md) or [CLIENT_INSTALL.md](CLIENT_INSTALL.md).
