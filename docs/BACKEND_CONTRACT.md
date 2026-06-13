# Backend Contract

Use Supabase Edge Functions in the Lovable project. The client should never talk directly to tables for sensitive writes. Use Edge Functions to validate, sign, and insert.

## Auth model

### Device linking

1. CLI calls `POST /device-start-link`
2. Backend creates one-time link code and URL.
3. CLI opens browser to `/me?link_code=...`
4. User signs in.
5. Website claims link code and attaches device to profile.
6. CLI polls `POST /device-complete-link` until linked.

## Endpoints

### POST /device-start-link

Request:

```json
{
  "device_public_key": "base64-ed25519-public-key",
  "device_name": "Kushagra-Windows",
  "client_version": "0.3.0",
  "build_id": "dev",
  "build_channel": "dev",
  "os": "windows",
  "arch": "amd64"
}
```

Response:

```json
{
  "device_id": "dev_...",
  "link_code": "abc123",
  "link_url": "https://sponsorbucks.xyz/me?link_code=abc123",
  "expires_at": "2026-06-13T12:00:00Z"
}
```

### POST /device-complete-link

Request:

```json
{
  "device_id": "dev_...",
  "link_code": "abc123"
}
```

Response if pending:

```json
{
  "status": "pending"
}
```

Response if linked:

```json
{
  "status": "linked",
  "device_token": "short-lived-or-refreshable-device-token",
  "user_id": "uuid"
}
```

### GET /ads-next

Query parameters:

- `surface`
- `country` optional
- `client_version`
- `last_creative_id` optional

Headers:

- `Authorization: Bearer <device_token>`

Response:

```json
{
  "campaign_id": "camp_...",
  "creative_id": "creative_...",
  "line": "Sponsored - Deploy APIs faster",
  "destination_url": "https://example.com",
  "display_ms": 5000,
  "rotation_allowed": true
}
```

The backend must enforce a 60-character maximum for `line`.

### POST /events-session-start

Signed JSON body:

```json
{
  "event_type": "session_start",
  "session_id": "sess_...",
  "device_id": "dev_...",
  "surface": "codex",
  "started_at": "2026-06-13T12:00:00Z",
  "client_version": "0.3.0",
  "build_id": "dev",
  "build_channel": "dev",
  "human_initiated": true
}
```

### POST /events-heartbeat

Signed JSON body:

```json
{
  "event_type": "heartbeat",
  "session_id": "sess_...",
  "device_id": "dev_...",
  "surface": "codex",
  "campaign_id": "camp_...",
  "creative_id": "creative_...",
  "creative_hash": "sha256-of-visible-line",
  "sequence": 4,
  "visible_ms": 5000,
  "screen_unlocked": true,
  "recent_input_bucket": "under_2_minutes",
  "foreground_supported_surface": true,
  "placement_visible": true,
  "created_at": "2026-06-13T12:00:05Z",
  "client_version": "0.3.0",
  "build_id": "dev",
  "build_channel": "dev"
}
```

### POST /events-session-end

Signed JSON body:

```json
{
  "event_type": "session_end",
  "session_id": "sess_...",
  "device_id": "dev_...",
  "surface": "codex",
  "ended_at": "2026-06-13T12:04:22Z",
  "exit_code": 0,
  "client_version": "0.3.0",
  "build_id": "dev",
  "build_channel": "dev"
}
```

## Signature headers

Every event endpoint should require:

```http
X-SponsorBucks-Device-Id: dev_...
X-SponsorBucks-Signature: base64(ed25519_signature_of_raw_body)
X-SponsorBucks-Timestamp: 2026-06-13T12:00:05Z
```

Backend rejects:

- invalid signature
- replayed sequence
- unknown device
- disabled device
- stale timestamp
- impossible concurrent sessions
