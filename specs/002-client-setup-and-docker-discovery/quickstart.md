# Quickstart: interactive setup and env-only Docker for bloop-tunnel

This guide shows the new operator paths for configuring `bloop-tunnel` without hand-authoring YAML unless you actually want to.

## Setup paths

### Option A: Interactive terminal setup (recommended)

Use the setup command to generate client configuration with guided prompts:

```bash
bloop-tunnel setup
```

Expected defaults:
- control plane URL: `https://api.bloop.to`
- relay URL: production hosted relay default
- auth token source: environment variable for safer local handling

Typical flow:
1. choose config output mode
2. accept or override hosted defaults
3. add one or more tunnels
4. optionally run Docker discovery
5. write config or print env-based output

## Example interactive output targets

### Write YAML config

```bash
bloop-tunnel setup --config ./client.yaml
```

### Print `.env` template

```bash
bloop-tunnel setup --output env-file > .env.bloop-tunnel
```

### Print Compose-ready environment block

```bash
bloop-tunnel setup --output compose-block
```

## Option B: Docker with environment variables only

Example Compose fragment:

```yaml
services:
  bloop-tunnel:
    image: ghcr.io/bloop/bloop-tunnel:latest
    environment:
      BLOOP_CONTROL_PLANE_URL: https://api.bloop.to
      BLOOP_RELAY_URL: wss://relay.bloop.to/connect
      BLOOP_AUTH_TOKEN_ENV: BLOOP_CLIENT_TOKEN
      BLOOP_CLIENT_TOKEN: replace-me

      BLOOP_TUNNELS_0_NAME: app
      BLOOP_TUNNELS_0_HOSTNAME: app.bloop.to
      BLOOP_TUNNELS_0_LOCAL_ADDR: host.docker.internal:3000
      BLOOP_TUNNELS_0_ACCESS: public

      BLOOP_TUNNELS_1_NAME: admin
      BLOOP_TUNNELS_1_LOCAL_ADDR: host.docker.internal:4000
      BLOOP_TUNNELS_1_ACCESS: basic_auth
      BLOOP_TUNNELS_1_BASIC_AUTH_USERNAME: gene
      BLOOP_TUNNELS_1_BASIC_AUTH_PASSWORD_ENV: BLOOP_ADMIN_PASSWORD
      BLOOP_ADMIN_PASSWORD: replace-me
    extra_hosts:
      - host.docker.internal:host-gateway
```

## Option C: Existing config-file workflow

The traditional config file path remains supported:

```bash
bloop-tunnel -config ./client.yaml
```

## Optional Docker discovery

If you want setup help for services already running in Docker, mount the Docker socket and request discovery explicitly.

Example local run:

```bash
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/bloop/bloop-tunnel:latest \
  bloop-tunnel setup --discover-docker
```

Behavior expectations:
- discovery is optional
- candidates are suggestions only
- no service is exposed automatically
- you must confirm each selected tunnel before it is written out

## Verification

After generating config or env variables, start the client and confirm:

- the client connects to the relay
- the configured tunnels register successfully
- hosted defaults point to `api.bloop.to` unless you intentionally overrode them

## Notes

- Use environment variables for secrets whenever possible.
- Keep local-development overrides explicit instead of mutating the production defaults in docs.
- If Docker discovery cannot access the daemon, continue with manual tunnel entry like a civilized person instead of failing the whole setup flow.
