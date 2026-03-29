# bloop-tunnel

`bloop-tunnel` is the production client package for exposing local HTTP services through `bloop.to`.

This repository’s release surface is **`bloop-client` only**: the agent that runs on a customer or operator machine and publishes configured local services.

## What `bloop-client` does

- enrolls with the control plane
- receives scoped runtime credentials
- connects outbound to the relay
- registers configured HTTP tunnels
- reports runtime status back to the control plane

## Setup paths

Use the path that fits the operator surface:

- **Interactive CLI**: `bloop-client setup --config ./client.yaml` walks through defaults, tunnel editing, and optional Docker discovery.
- **Config file**: pass `--config <path>` to run from YAML.
- **Env-only Docker**: omit `--config` and define zero-based indexed `BLOOP_TUNNELS_<n>_*` entries directly (for example `BLOOP_TUNNELS_0_*`, `BLOOP_TUNNELS_1_*`).
- **Automation scaffold**: `bloop-client setup --non-interactive` still generates starter YAML, `.env`, or Compose-friendly output.

Hosted defaults point to:

- control plane: `https://api.bloop.to`
- relay: `wss://relay.bloop.to/connect`

## Install and run

See [`docs/CLIENT_INSTALL.md`](docs/CLIENT_INSTALL.md) for:

- release artifact names by platform
- Docker usage
- native macOS, Linux, and Windows usage
- config shape and examples
- verification steps

## Release expectations

GitHub releases for this repository publish native archives for `bloop-client` on:

- Linux (`amd64`, `arm64`)
- macOS (`amd64`, `arm64`)
- Windows (`amd64`)

Expected asset names:

- `bloop-client-vX.Y.Z-linux-amd64.tar.gz`
- `bloop-client-vX.Y.Z-linux-arm64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-amd64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-arm64.tar.gz`
- `bloop-client-vX.Y.Z-windows-amd64.zip`

Each archive contains:

- `bloop-client`
- the client install guide
- example client configuration

## Local build

Build the client locally:

```bash
go build -o dist/bloop-client ./cmd/bloop-client
```

Example cross-build:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bloop-client-windows-amd64.exe ./cmd/bloop-client
```

## Scope

This repository is the canonical home for the end-user and operator machine-local client runtime: `bloop-client`.
