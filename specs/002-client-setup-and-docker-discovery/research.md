# Research Notes: bloop-client interactive setup, production defaults, and Docker endpoint discovery

## Chosen Direction

- Add a dedicated interactive setup path to the client instead of forcing manual YAML editing.
- Keep existing file-based configuration working for backward compatibility.
- Add first-class environment-variable parsing for unlimited tunnel definitions using indexed keys.
- Make `https://api.bloop.to` the default control-plane host for generated config and docs.
- Keep relay defaults production-oriented but overridable for local development.
- Make Docker socket discovery opt-in, local-only in spirit, and confirmation-based.
- Treat docs and frontend alignment as part of the feature, not afterthought cleanup.

## Alternatives Considered

### Keep config-file-only setup
Rejected because it preserves the main operator pain point and makes Docker usage more awkward than it needs to be.

### Use a single JSON blob environment variable for all tunnels
Rejected as the primary path because it is harder to read, harder to maintain in Compose, and easier to break silently. It can remain a future enhancement if needed.

### Make env-defined tunnels replace the entire file config implicitly
Rejected for now because silent replacement is too surprising. Precedence and merge behavior should be explicit and tested.

### Auto-expose discovered Docker services immediately
Rejected because it is reckless, confusing, and exactly how you end up publishing something dumb to the internet by accident.

### Build Docker discovery into the runtime startup path
Rejected because discovery belongs in setup guidance, not in every normal client start.

### Require frontend changes before backend/client work can ship
Rejected because docs and frontend should align quickly, but the client setup improvements should not be blocked on a separate repo if the runtime behavior is ready.
