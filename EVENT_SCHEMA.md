# Event Schema

## session_start

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

## heartbeat

```json
{
  "event_type": "heartbeat",
  "session_id": "sess_...",
  "device_id": "dev_...",
  "surface": "codex",
  "campaign_id": "camp_...",
  "creative_id": "creative_...",
  "creative_hash": "sha256-of-visible-line",
  "sequence": 1,
  "visible_ms": 5000,
  "screen_unlocked": true,
  "recent_input_bucket": "unknown_v0",
  "foreground_supported_surface": true,
  "placement_visible": true,
  "created_at": "2026-06-13T12:00:05Z",
  "client_version": "0.3.0",
  "build_id": "dev",
  "build_channel": "dev"
}
```

## session_end

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
