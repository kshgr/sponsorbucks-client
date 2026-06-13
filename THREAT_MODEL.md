# Threat Model

## Primary risks

- fake sessions
- replayed heartbeats
- hidden or unattended agents
- disabled or paused placements being bypassed
- privacy regression through new telemetry

## Defenses

- signed event bodies
- session IDs and heartbeat sequence numbers
- local pause and per-tool disable controls
- localhost-only daemon access
- no workspace scanning
- no command-argument logging after `--` by default

## Out of scope for v0

- auto-update
- remote executable download
- code or terminal capture

