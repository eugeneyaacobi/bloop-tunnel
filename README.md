# bloop-tunnel

Private-first HTTP tunnel client for exposing local services through `bloop.to`.

This repository is the **client-side package**: it is for the user/local machine that connects outward to a relay. The relay/server component should be packaged with the rest of the server infrastructure or moved to its own repository.

## Status

Early implementation. See `specs/001-v1-http-tunnels/` for the original combined-system spec and implementation notes.

## What this repo ships

This repo's public install and release story is centered on one binary:

- `bloop-client` — the local agent that connects to the relay and exposes your local HTTP service

`bloop-relay` is server-side infrastructure and is no longer part of this repository's user-facing native release artifacts.

## Runtime ingest (v1)

The production-shaped ingest path is client-owned:
- the client enrolls with control-plane
- the client receives a scoped ingest token
- the client reports runtime truth directly to control-plane

## Install the client

See `docs/CLIENT_INSTALL.md` for:
- Docker install/run
- native macOS/Linux/Windows usage
- config examples
- verification steps
- AI agent / automation hints


## CI and release artifacts

GitHub Actions handles verification and client-native release packaging:

- `.github/workflows/ci.yml`
  - runs `go test ./...`
  - cross-builds `bloop-client` for Linux, macOS, and Windows
  - uploads per-platform client artifacts for PR/push validation

- `.github/workflows/release.yml`
  - triggers on `v*` tags, published releases, or manual dispatch
  - builds versioned native archives for `bloop-client`
  - uploads release assets to the GitHub Release

### Expected GitHub release assets

Each release should include archives like:

- `bloop-client-vX.Y.Z-linux-amd64.tar.gz`
- `bloop-client-vX.Y.Z-linux-arm64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-amd64.tar.gz`
- `bloop-client-vX.Y.Z-darwin-arm64.tar.gz`
- `bloop-client-vX.Y.Z-windows-amd64.zip`

Each archive contains:
- `bloop-client`
- client install docs
- example client config

## Local build examples

Build the client locally:

```bash
go build -o dist/bloop-client ./cmd/bloop-client
```

Cross-build an example native artifact locally:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bloop-client-windows-amd64.exe ./cmd/bloop-client
```

## Server-side relay

The server-side relay now belongs in the separate `bloop-relay` repository/package. This repository is focused on `bloop-client` only.
