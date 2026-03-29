# Control Plane API Contract

**API Version**: v1  
**Base URL**: `{control_plane_url}`  
**Purpose**: Enrollment and runtime registration for bloop-tunnel

---

## Overview

The control plane API handles client enrollment, issuing installation IDs and ingest tokens for runtime telemetry. The TUI verifies connectivity and enrollment tokens against this API.

---

## Endpoints

### 1. Health Check

Used by TUI to verify control plane connectivity.

```
GET /health
```

**Request**:
- Headers: None required
- Body: None

**Response (200 OK)**:
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

**Response (503 Service Unavailable)**:
```json
{
  "status": "unhealthy",
  "error": "Database connection failed"
}
```

**TUI Behavior**:
- On success: Proceed to enrollment verification
- On failure: Show error with specific details (DNS, TLS, HTTP)
- Offer retry, skip, or exit options

---

### 2. Enroll Runtime

Used by client to register a runtime installation and receive credentials for telemetry ingest.

```
POST /api/runtime/enroll
```

**Request**:
- Headers:
  - `Content-Type: application/json`
- Body:
```json
{
  "enrollment_token": "enroll_abc123xyz...",
  "runtime_name": "default-client"
}
```

**Fields**:
- `enrollment_token` (string, required): Token issued by control plane for enrollment
- `runtime_name` (string, required): Human-readable name for this runtime instance

**Response (201 Created)**:
```json
{
  "installation_id": "install_abc123xyz...",
  "ingest_token": "ingest_xyz789uvw...",
  "ingest_endpoint": "wss://ingest.bloop.to/ingest"
}
```

**Fields**:
- `installation_id` (string): Unique identifier for this runtime installation
- `ingest_token` (string): Bearer token for authenticated ingest WebSocket connection
- `ingest_endpoint` (string): WebSocket endpoint for telemetry ingest

**Response Errors**:

**400 Bad Request**:
```json
{
  "error": "invalid_request",
  "message": "enrollment_token is required"
}
```

**401 Unauthorized**:
```json
{
  "error": "invalid_token",
  "message": "The provided enrollment token is invalid or expired"
}
```

**409 Conflict**:
```json
{
  "error": "already_enrolled",
  "message": "This installation is already enrolled with a different token"
}
```

**429 Too Many Requests**:
```json
{
  "error": "rate_limited",
  "message": "Too many enrollment attempts. Please try again later.",
  "retry_after": 60
}
```

**500 Internal Server Error**:
```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

**TUI Behavior**:
- On success (201): Display installation ID, proceed to next step
- On 401 invalid_token: Show "Invalid or expired enrollment token" error, offer retry
- On 401 unauthorized: Show "Unauthorized enrollment token" error, offer retry
- On 409: Show "Already enrolled" error, recommend using existing installation
- On 429: Show rate limit error with retry-after time, offer wait or skip
- On 5xx: Show "Control plane error" message, offer retry or skip
- Timeout (10s): Show "Request timed out" error, offer retry

---

## Connectivity Verification

The TUI performs connectivity verification in three stages:

### Stage 1: DNS Resolution

Resolve the hostname from `ControlPlaneURL` to IP address.

**Success**: DNS resolved, IP address available
**Failure**: DNS resolution failed (NXDOMAIN, timeout, network error)

**TUI Display**:
- Success: "DNS resolved to {IP}"
- Failure: "DNS resolution failed: {error}"

### Stage 2: TLS Handshake

Perform TLS handshake if using HTTPS.

**Success**: TLS handshake completed, certificate valid
**Failure**: TLS handshake failed (certificate expired, invalid, network error)

**TUI Display**:
- Success: "TLS connection established"
- Failure: "TLS handshake failed: {error}"

### Stage 3: HTTP Health Check

Send GET request to `/health` endpoint.

**Success**: 200 OK response received
**Failure**: Non-200 response, timeout, network error

**TUI Display**:
- Success: "Control plane healthy (latency: {latency})"
- Failure: "Health check failed: {error}"

---

## Security Considerations

### Token Handling

- Enrollment tokens are long-lived secrets issued by the control plane
- Tokens must be stored in environment variables, never in plaintext in config files
- Ingest tokens are returned from enrollment and should be stored in environment variables
- Tokens must never be logged or displayed in the TUI

### API Security

- All API requests use HTTPS in production
- TLS certificates must be validated
- Enrollment tokens are single-use or have limited validity
- Rate limiting prevents brute force token guessing

### TUI Security

- Mask enrollment token input using password field type
- Never display tokens in plaintext in TUI views
- Clear token values from memory after enrollment
- Redact tokens from error messages and logs

---

## Error Codes

| Error Code | HTTP Status | Description | TUI Action |
|------------|-------------|-------------|------------|
| `invalid_request` | 400 | Missing or malformed request | Show specific field error |
| `invalid_token` | 401 | Token is invalid | Show "Invalid token" error |
| `expired_token` | 401 | Token has expired | Show "Expired token" error |
| `unauthorized` | 401 | Token not authorized for enrollment | Show "Unauthorized" error |
| `already_enrolled` | 409 | Installation already enrolled | Show warning, recommend existing |
| `rate_limited` | 429 | Too many requests | Show retry-after time |
| `internal_error` | 500 | Server error | Show generic error, offer retry |

---

## Timeouts

The TUI uses the following timeouts for control plane operations:

- DNS Resolution: 5s
- TLS Handshake: 5s
- HTTP Health Check: 5s
- Enrollment Request: 10s

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

## Testing

### Unit Testing

Mock HTTP responses for:
- Success (200, 201)
- Client errors (400, 401)
- Server errors (500)
- Timeout scenarios

### Integration Testing

Test against real control plane endpoint:
- Valid enrollment token
- Invalid enrollment token
- Expired enrollment token
- Network error scenarios

---

## Notes

- The control plane URL defaults to `https://api.bloop.to` in production
- Health check endpoint is used for connectivity verification only
- Enrollment token verification is optional but recommended
- Installation ID and ingest token are set as environment variables for runtime
- The TUI stores the installation ID in config but not the ingest token (use env var)
