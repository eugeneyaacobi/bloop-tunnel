# bloop-tunnel

`bloop-tunnel` is the production client package for exposing local HTTP services through `bloop.to`.

This repository’s release surface is **`bloop-client` only**: the agent that runs on a customer or operator machine, connects outbound to the relay, and publishes configured local services.

Server-side relay infrastructure is not part of this package’s install or release story.

## What `bloop-client` does

- enrolls with the control plane
- receives scoped runtime credentials
- connects outbound to the relay
- registers configured HTTP tunnels
- reports runtime status back to the control plane

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
