# Frontend Implementation Checklist: client setup + Docker discovery alignment

Purpose: turn the spec's frontend-alignment work into a shippable execution checklist the frontend repo can implement without re-reading the whole spec stack.

## Scope

Align frontend setup guidance, install surfaces, and runtime terminology with spec 002:
- terminal-first setup flow
- env-only Docker configuration
- production defaults pointing to `https://api.bloop.to`
- opt-in Docker discovery with explicit confirmation
- config-file workflow kept as fallback, not the only path

## Priority order

1. **P1 — Fix misleading setup model**
2. **P1 — Fix production-default copy**
3. **P2 — Add env-only Docker guidance**
4. **P2 — Add optional Docker discovery explanation**
5. **P2 — Align runtime/status language with the new model**

---

## Execution checklist

### 1) Audit and replace stale setup assumptions

- [ ] Find all frontend strings and screens that imply users must hand-edit a config file.
- [ ] Replace config-file-only framing with a three-path model:
  - terminal setup
  - Docker env setup
  - manual config file fallback
- [ ] Remove or rewrite any setup text that treats localhost defaults as the normal hosted path.
- [ ] Flag any UI copy that implies Docker discovery is automatic exposure.

**Candidate surfaces**
- onboarding/setup wizard
- install modal / client install drawer
- docs panel embedded in app
- tunnel/connector empty states
- client-agent status/help panels
- “how to connect” screens

**Acceptance check**
- A user can identify the recommended setup path without seeing contradictory instructions.

---

### 2) Update setup-path information architecture

- [ ] Present setup choices in this order:
  1. **Terminal setup** — recommended for most users
  2. **Docker env setup** — recommended for Compose/container users
  3. **Manual config file** — fallback for advanced/manual workflows
- [ ] Add brief tradeoff text for each path.
- [ ] Ensure UI uses consistent nouns across surfaces:
  - use **terminal setup** instead of vague “CLI mode”
  - use **environment variables** or **env-only Docker setup** instead of “advanced Docker config”
  - use **Docker discovery** instead of “auto-expose” or “automatic publish”

**Component implications**
- likely one shared “Setup Paths” content block or card group
- avoid duplicating path descriptions across multiple screens; centralize copy source if possible

**Acceptance check**
- All setup surfaces present the same 3-path model and same naming.

---

### 3) Correct production defaults copy

- [ ] Update setup guidance to show `https://api.bloop.to` as the default control-plane endpoint.
- [ ] Remove any localhost-based hosted defaults from examples, helper text, placeholders, or screenshots.
- [ ] If relay/default endpoint text appears in UI, ensure it matches the backend/docs decision rather than old local-dev assumptions.
- [ ] Make local-development overrides explicit as overrides, not defaults.

**Copy changes**
- Prefer: “Defaults target the hosted bloop.to control plane.”
- Prefer: “Override endpoints only for local development or custom deployments.”
- Avoid: “Change localhost to your production URL when ready.”

**Acceptance check**
- No primary setup surface implies localhost is the normal hosted starting point.

---

### 4) Add env-only Docker setup guidance

- [ ] Explain that Docker users can configure tunnels entirely through environment variables.
- [ ] Show env-only setup as first-class, not as a workaround.
- [ ] Reference support for multiple indexed tunnel definitions.
- [ ] Ensure examples and help text do not imply a mounted YAML file is required.

**Candidate surfaces**
- Docker install tab
- Compose example drawer
- deployment method switcher
- client install docs embedded in frontend

**Copy changes**
- Prefer: “Run bloop-tunnel in Docker using environment variables only.”
- Prefer: “Define one or more tunnels with indexed env vars.”
- Avoid: “Mount a config file, or optionally pass some env overrides.”

**Acceptance check**
- A Docker-first user can understand that fileless setup is supported.

---

### 5) Explain Docker discovery safely

- [ ] Describe Docker discovery as **optional**.
- [ ] State clearly that discovery only suggests running services; it does not expose anything by itself.
- [ ] State clearly that users must explicitly confirm each discovered service before it becomes a tunnel.
- [ ] If UI includes feature badges/tooltips, label discovery as a local convenience feature.
- [ ] Add error/help copy for unavailable Docker socket cases if surfaced in frontend guidance.

**Candidate surfaces**
- setup tips/help text
- Docker-specific onboarding
- feature comparison table
- FAQ / troubleshooting panel

**Copy changes**
- Prefer: “Optionally inspect local Docker services during setup.”
- Prefer: “Discovered services are suggestions until you confirm them.”
- Avoid: “Automatically expose your Docker services.”

**Acceptance check**
- Frontend copy cannot be reasonably read as promising automatic publication.

---

### 6) Align runtime/status and empty-state messaging

- [ ] Update empty states to reference the new ways a client can be configured.
- [ ] If runtime status/help screens mention expected config sources, include file and env-based setups.
- [ ] Ensure tunnel-management help does not assume tunnels were manually authored in YAML.
- [ ] If the frontend surfaces setup completion or next steps, ensure it points users toward terminal setup and Docker env setup appropriately.

**Candidate surfaces**
- no-client-connected empty state
- no-tunnels-configured empty state
- client detail/status page
- install success / next-step callouts

**Acceptance check**
- Runtime-oriented surfaces describe config origin neutrally and do not bias toward legacy YAML-only mental models.

---

## Suggested shared copy blocks

### Setup paths summary

> Choose the setup path that matches how you run bloop-tunnel: use terminal setup for the guided default flow, env-only Docker setup for container deployments, or a manual config file if you want direct file control.

### Production defaults summary

> New setups default to the hosted bloop.to control plane at `https://api.bloop.to`. Only override endpoints for local development or custom deployments.

### Docker discovery summary

> Docker discovery is optional. It can suggest running local services during setup, but nothing is exposed until you explicitly confirm each tunnel.

---

## Acceptance criteria for the frontend PR

- [ ] Terminal setup is presented as the recommended default path.
- [ ] Docker env-only setup is presented as a first-class path.
- [ ] Manual config-file setup is still documented, but framed as fallback/manual control.
- [ ] `https://api.bloop.to` is the visible default on primary setup surfaces.
- [ ] No frontend surface implies Docker discovery auto-exposes services.
- [ ] Setup terminology is consistent across onboarding, install, help, and status surfaces.
- [ ] Any screenshots/examples/placeholders used by the frontend match the new defaults.

---

## Dependencies

### Upstream dependencies
- Backend/docs must confirm the final hosted relay default text if surfaced in UI.
- CLI setup command names/flags should be stable enough to reference in frontend guidance.
- Docker env variable naming scheme should be stable enough for examples/help text.

### Frontend implementation dependencies
- Shared copy source or centralized content config is strongly preferred to prevent setup-copy drift.
- If onboarding/install content is split across multiple components, update all of them in the same PR.

### Coordination risks
- If the frontend ships before backend/docs naming stabilizes, copy drift is likely around:
  - command names
  - relay endpoint wording
  - Docker discovery phrasing
  - env var example structure

---

## Recommended implementation slices

### Slice A — Copy correction pass (ship first)
- remove localhost defaults
- remove YAML-only framing
- remove automatic-discovery wording

### Slice B — Setup path UI alignment
- add 3-path setup model
- add short comparison/tradeoff content
- update empty states and install surfaces

### Slice C — Docker-specific polish
- add env-only Docker examples/help
- add Docker discovery explainer and caveats
- review tooltips/FAQ/troubleshooting

---

## Review checklist

- [ ] Product review: recommended path is obvious
- [ ] Frontend review: no duplicated conflicting copy remains
- [ ] Accessibility review: path selection and help text remain readable and scannable
- [ ] Security/trust review: Docker discovery wording does not overpromise or hide confirmation requirements
- [ ] Spec review: matches FR-017, FR-018, SC-006
