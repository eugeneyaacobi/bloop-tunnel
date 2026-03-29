# Relay API Contract

**API Version**: v1  
**Protocol**: WebSocket  
**Relay URL**: `{relay_url}`  
**Purpose**: Tunnel registration and WebSocket connection for proxying traffic

---

## Overview

The relay is a WebSocket-based service that accepts tunnel registrations from clients and forwards incoming connections to local services. The TUI verifies relay connectivity and auth token validity against this service.

---

## Connection Protocol

### 1. Establish WebSocket Connection

Connect to the relay WebSocket endpoint.

```
WS {relay_url}/connect
```

**Connection Headers**:
- `Authorization: Bearer {auth_token}`
- `X-Client-Version: {client_version}`

**Query Parameters**:
- `installation_id`: Installation ID from enrollment (optional, for telemetry)

**Connection Flow**:
1. Establish WebSocket connection with auth token
2. Receive handshake confirmation from relay
3. Send `RegisterTunnels` message to register tunnels
4. Receive `RegisterTunnelsAck` with tunnel URLs
5. Begin proxying traffic

### 2. Handshake Confirmation

After successful WebSocket connection, relay sends confirmation:

```json
{
  "type": "connected",
  "timestamp": "2026-03-29T19:30:00Z",
  "relay_version": "1.0.0"
}
```

**Fields**:
- `type` (string): Message type identifier
- `timestamp` (string): ISO 8601 timestamp
- `relay_version` (string): Relay server version

### 3. Register Tunnels

Client sends tunnel registration message:

```json
{
  "type": "RegisterTunnels",
  "tunnels": [
    {
      "name": "app",
      "local_addr": "127.0.0.1:3000",
      "access": "public"
    },
    {
      "name": "admin",
      "local_addr": "127.0.0.1:4000",
      "access": "basic_auth",
      "basic_auth": {
        "username": "admin",
        "password_env": "ADMIN_PASSWORD"
      }
    }
  ]
}
```

**Tunnel Fields**:
- `name` (string, required): Unique tunnel identifier
- `local_addr` (string, required): Local address to forward to (host:port)
- `access` (string, required): Access mode - "public", "basic_auth", "token_protected"
- `basic_auth` (object, optional): Basic auth configuration (required if access = "basic_auth")
  - `username` (string, required): Basic auth username
  - `password_env` (string, required): Environment variable containing password
- `token_env` (string, optional): Environment variable containing access token (required if access = "token_protected")
- `hostname` (string, optional): Custom hostname for tunnel (default: `{name}.{domain}`)

### 4. Register Tunnels Acknowledgment

Relay responds with tunnel URLs:

```json
{
  "type": "RegisterTunnelsAck",
  "tunnels": [
    {
      "name": "app",
      "url": "https://app.bloop.to",
      "status": "registered"
    },
    {
      "name": "admin",
      "url": "https://admin.bloop.to",
      "status": "registered"
    }
  ]
}
```

**Fields**:
- `type` (string): Message type identifier
- `tunnels` (array): Array of tunnel registrations
  - `name` (string): Tunnel name (matches request)
  - `url` (string): Public URL for the tunnel
  - `status` (string): "registered" or "error"

### 5. Error Messages

Relay sends error messages for failures:

```json
{
  "type": "error",
  "code": "unauthorized",
  "message": "Invalid or expired auth token"
}
```

**Error Codes**:
- `unauthorized`: Invalid or expired auth token
- `invalid_request`: Malformed request
- `tunnel_exists`: Tunnel name already registered
- `tunnel_invalid`: Invalid tunnel configuration
- `rate_limited`: Too many registration attempts
- `internal_error`: Relay server error

---

## Relay Token Verification

The TUI verifies relay auth token by:

1. Establishing WebSocket connection with auth token
2. Waiting for handshake confirmation (`connected` message)
3. Closing connection after confirmation

**Success Criteria**:
- WebSocket connection established successfully
- `connected` message received from relay
- No `error` messages received

**Failure Scenarios**:
- Connection refused (401): Invalid or expired token
- Connection timeout: Relay unreachable
- Network error: Connectivity issues
- `error` message with `unauthorized`: Token invalid
- `error` message with other code: Relay error

---

## Connection States

```
[Disconnected]
      |
      v
[Connecting] - WebSocket handshake in progress
      |
      v (success)
[Connected] - Handshake confirmation received
      |
      v (send RegisterTunnels)
[Registering] - Waiting for tunnel registration
      |
      v (success)
[Active] - Tunnels registered, proxying traffic
```

**Error States**:
- `[ConnectionFailed]` - Could not establish connection
- `[Unauthorized]` - Auth token invalid
- `[RegistrationFailed]` - Tunnel registration failed

---

## Security Considerations

### Auth Token Handling

- Auth tokens are long-lived secrets issued by the relay service
- Tokens must be stored in environment variables, never in plaintext in config files
- Tokens are sent via WebSocket `Authorization` header
- Tokens must never be logged or displayed in the TUI

### WebSocket Security

- All WebSocket connections use `wss://` (secure WebSocket) in production
- TLS certificates must be validated
- Auth tokens are validated on connection
- Rate limiting prevents brute force attacks

### TUI Security

- Mask relay auth token input using password field type
- Never display tokens in plaintext in TUI views
- Clear token values from memory after verification
- Redact tokens from error messages and logs

---

## Timeouts

The TUI uses the following timeouts for relay operations:

- WebSocket Connection: 10s
- Handshake Confirmation: 5s
- Registration Acknowledgment: 5s

If any operation times out, the TUI shows a timeout error and offers retry options.

---

## Retry Logic

The TUI implements exponential backoff for retries:

- Retry 1: Immediate
- Retry 2: Wait 1s
- Retry 3: Wait 2s
- Retry 4: Wait 4s
- Retry 5: Wait 8s

Maximum retries: 5
Total max wait time: 15s

If all retries fail, show error and offer skip or exit options.

---

## Error Codes

| Error Code | Description | TUI Action |
|------------|-------------|------------|
| `unauthorized` | Invalid or expired auth token | Show "Invalid relay token" error |
| `invalid_request` | Malformed request | Show "Invalid request" error |
| `tunnel_exists` | Tunnel name already registered | Show "Tunnel exists" warning |
| `tunnel_invalid` | Invalid tunnel configuration | Show "Invalid tunnel" error |
| `rate_limited` | Too many attempts | Show rate limit error with retry-after |
| `internal_error` | Relay server error | Show generic error, offer retry |
| `connection_failed` | Could not establish connection | Show "Connection failed" error |
| `timeout` | Operation timed out | Show "Timeout" error |

---

## Connectivity Verification

The TUI performs relay connectivity verification in two stages:

### Stage 1: WebSocket Connection

Establish WebSocket connection to `{relay_url}/connect` with auth token.

**Success**: WebSocket connection established, handshake completed
**Failure**: Connection refused, timeout, network error

**TUI Display**:
- Success: "Connected to relay"
- Failure: "Connection failed: {error}"

### Stage 2: Handshake Confirmation

Wait for `connected` message from relay.

**Success**: `connected` message received
**Failure**: `error` message received or timeout

**TUI Display**:
- Success: "Relay handshake successful"
- Failure: "Handshake failed: {error}"

---

## Testing

### Unit Testing

Mock WebSocket server for:
- Successful connection and handshake
- Unauthorized connection (401)
- Connection timeout
- Server error scenarios

### Integration Testing

Test against real relay endpoint:
- Valid auth token
- Invalid auth token
- Expired auth token
- Network error scenarios

---

## Notes

- The relay URL defaults to production relay endpoint in production
- WebSocket connections use `wss://` in production
- Auth token verification is optional but recommended
- Relay tokens are long-lived and should be rotated periodically
- The TUI verifies connectivity but does not complete full registration during setup
- Full tunnel registration happens on client startup (`bloop-client run`)
