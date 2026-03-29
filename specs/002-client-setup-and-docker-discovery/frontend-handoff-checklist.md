# Frontend Handoff Checklist: spec 002

This artifact tightens the frontend follow-through for spec 002 so Prism/front-end implementation can execute without guessing.

## Scope boundary

Frontend does **not** implement Docker discovery itself.
Frontend **does** explain the setup modes, production defaults, and the fact that discovery is optional, local, and confirmation-based.

## Required UI copy changes

- [ ] Replace any config-file-only setup copy with a three-path model:
  - interactive CLI setup
  - env-only Docker setup
  - YAML config fallback
- [ ] Default hosted endpoint references to `https://api.bloop.to`
- [ ] Use `wss://relay.bloop.to/connect` wherever relay examples are shown
- [ ] Describe Docker discovery as:
  - opt-in
  - local-only in spirit
  - dependent on mounted Docker socket
  - requiring explicit confirmation before any tunnel is added
- [ ] Remove any implication that discovery auto-publishes containers

## Setup surfaces to update

- [ ] install/onboarding page
- [ ] runtime/installation detail guidance
- [ ] empty states for "no client configured"
- [ ] any generated command snippets or copyable env templates
- [ ] help/FAQ surfaces that currently mention YAML-only setup

## Content requirements

- [ ] Show `bloop-tunnel setup` as the primary guided path
- [ ] Show `bloop-tunnel setup --output env-file` for Docker operators
- [ ] Explain env-only tunnel encoding with `BLOOP_TUNNELS_<n>_*`
- [ ] Call out precedence in plain language: flags > env > config file > defaults
- [ ] Mention that tunnel definitions discovered from Docker become explicit saved config, not hidden runtime state

## Acceptance checklist

- [ ] A new user can tell which setup path to use without reading backend docs first
- [ ] A Docker operator can copy an env-only example without needing YAML
- [ ] Discovery copy never implies automatic exposure
- [ ] Hosted defaults shown in UI match repo docs exactly
- [ ] Terminology matches CLI/docs: "setup", "env-only", "Docker discovery", "control plane", "relay"

## Handoff note for Prism

Treat this as copy-and-flow alignment work, not a new orchestration feature. The backend owns discovery execution. The frontend owns accurate explanation and operator expectations.
