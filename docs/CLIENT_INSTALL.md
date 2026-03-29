# bloop-client install guide

This guide covers installing and running `bloop-client`, the machine-local agent for exposing HTTP services through `bloop.to`.

## What `bloop-client` does

`bloop-client`:

- enrolls with the control plane using an enrollment token
- receives scoped runtime credentials
- connects outbound to the relay
- registers configured tunnels
- reports runtime status to the control plane

`bloop-relay` is server infrastructure and is not part of this package’s normal install flow.

## Hosted defaults

Unless you explicitly override them for local testing, generated config and docs assume:

- control plane: `https://api.bloop.to`
- relay: `wss://relay.bloop.to/connect`

## Release artifacts

GitHub releases include native archives for:

- Linux (`amd64`, `arm64`)
- macOS (`amd64`, `arm64`)
- Windows (`amd64`)

Archive naming:

- `bloop-client-vX.Y.Z-linux-amd64.tar.gz`
- `bloop-client-vX.Y.Z-linux-arm64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-amd64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-arm64.tar.gz`
- `bloop-client-vX.Y.Z-windows-amd64.zip`

Each archive contains:

- `bloop-client`
- this install guide
- example client configuration

## Option A: Interactive terminal setup

Run the guided setup flow:

```bash
./bloop-client setup --config ./client.yaml
```

The setup flow now walks you through:

- control plane and relay defaults
- token environment variable names
- reconnect settings
- add/edit/remove tunnel definitions
- optional Docker discovery when explicitly requested

If `./client.yaml` already exists, setup preloads those values as prompt defaults so you can edit instead of starting over.

For automation or CI, you can still skip prompts entirely:

```bash
./bloop-client setup --non-interactive --config ./client.yaml
./bloop-client setup --non-interactive --output env-file > .env.bloop-client
./bloop-client setup --non-interactive --output compose-block
```

## Option B: Docker with env-only configuration

Build locally:

```bash
docker build -f deploy/docker/client.Dockerfile -t bloop-client:local .
```

Run with only environment variables:

```bash
docker run --rm \
  --env-file ./deploy/examples/client.env.example \
  -e BLOOP_CLIENT_TOKEN=replace-me \
  -e BLOOP_ENROLLMENT_TOKEN=optional-enrollment-token \
  --add-host=host.docker.internal:host-gateway \
  bloop-client:local
```

The client scans all zero-based indexed tunnel variables matching `BLOOP_TUNNELS_<n>_*` (for example `BLOOP_TUNNELS_0_*`, `BLOOP_TUNNELS_1_*`), so Docker and Compose deployments can define any number of tunnels.

Supported per-tunnel env fields include:

- `NAME`
- `HOSTNAME`
- `LOCAL_ADDR`
- `ACCESS`
- `TOKEN`
- `TOKEN_ENV`
- `BASIC_AUTH_USERNAME`
- `BASIC_AUTH_PASSWORD`
- `BASIC_AUTH_PASSWORD_ENV`

Precedence model:

1. built-in production defaults
2. config file values, when `--config` is provided
3. environment variables, including indexed tunnel definitions

If any indexed tunnel env entries are present, the environment-defined tunnel set replaces file-defined tunnels.

## Optional Docker discovery during setup

If your apps already run in Docker and you want help turning them into tunnel definitions, request discovery explicitly:

```bash
./bloop-client setup --config ./client.yaml --discover-docker
```

Notes:

- discovery is optional
- discovery only inspects the local Docker daemon when the socket is available
- setup asks for confirmation before adding each discovered service as a tunnel
- if Docker or the socket is unavailable, setup warns and continues with manual tunnel entry

## Option C: Native binary with YAML config

### Download a release artifact

Download the archive for your OS and architecture from the project’s GitHub Release, then extract it.

### Or build locally

```bash
go build -o bloop-client ./cmd/bloop-tunnel
```

### macOS / Linux

```bash
./bloop-client --config ./deploy/examples/client.example.yaml
```

### Windows (PowerShell)

```powershell
.\bloop-client.exe --config .\deploy\examples\client.example.yaml
```

## Example YAML config

```yaml
control_plane_url: https://api.bloop.to
relay_url: wss://relay.bloop.to/connect
auth_token_env: BLOOP_CLIENT_TOKEN
enrollment_token_env: BLOOP_ENROLLMENT_TOKEN
reconnect:
  initial_delay_ms: 1000
  max_delay_ms: 30000
logging:
  level: info
  format: json
tunnels:
  - name: app
    hostname: app.bloop.to
    local_addr: host.docker.internal:3000
    access: public
```

## Verify the client is healthy

Look for these client log signals:

- `runtime enrollment succeeded`
- `connected to relay`
- `client registered tunnels successfully`

Then verify in the control plane:

- installation status is active
- last seen is updating
- tunnel runtime reflects the configured routes

## Automation flow

For scripted installs or agent-driven setup:

1. generate starter YAML or env scaffolding with `bloop-client setup --non-interactive`
2. provide `BLOOP_CLIENT_TOKEN` and optionally `BLOOP_ENROLLMENT_TOKEN`
3. define one or more tunnels via YAML or zero-based indexed `BLOOP_TUNNELS_<n>_*` variables
4. run `bloop-client` with `--config <file>` or env-only Docker startup
5. confirm the installation becomes active and routes appear in the control plane
