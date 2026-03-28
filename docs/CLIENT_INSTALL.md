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

## Option A: Docker

Build locally:

```bash
docker build -f deploy/docker/client.Dockerfile -t bloop-client:local .
```

Run with a config file:

```bash
docker run --rm \
  -v $(pwd)/deploy/examples/client.relay-ingest.yaml:/config/client.yaml:ro \
  --add-host=host.docker.internal:host-gateway \
  bloop-client:local /bloop-client -config /config/client.yaml
```

## Option B: Native binary

### Download a release artifact

Download the archive for your OS and architecture from the project’s GitHub Release, then extract it.

### Or build locally

```bash
go build -o bloop-client ./cmd/bloop-client
```

### macOS / Linux

```bash
./bloop-client -config ./deploy/examples/client.relay-ingest.yaml
```

### Windows (PowerShell)

```powershell
.\bloop-client.exe -config .\deploy\examples\client.relay-ingest.yaml
```

## Required config values

Example configuration:

```yaml
relay_url: ws://relay.example/connect
auth_token: dev-token
control_plane_url: https://api.example.com
enrollment_token: bloop_enr_xxx
reconnect:
  initial_delay_ms: 1000
  max_delay_ms: 5000
logging:
  level: info
  format: json
tunnels:
  - name: app
    hostname: app.example.com
    local_addr: 127.0.0.1:8080
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

1. create an installation in the control plane
2. capture the enrollment token
3. write a client config with `control_plane_url`, `relay_url`, `auth_token`, `enrollment_token`, and the tunnel list
4. run `bloop-client -config <file>`
5. confirm the installation becomes active and routes appear in the control plane
