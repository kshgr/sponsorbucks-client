# SponsorBucks Client

SponsorBucks is an open-source local client for supported AI agent tools. The public repo is https://github.com/kshgr/sponsorbucks-client.

This technical preview is intentionally narrow:

- No code collection.
- No prompt collection.
- No terminal-output collection.
- Ads are text plus URL only.
- No remote executable ad payloads.
- No auto-update in v0.
- Events are signed locally before upload.
- Source and dev builds are limited until official signing is active.


## Current launch status

This is a technical preview for waitlist launch. Demo sessions work locally. Paid campaigns and payouts are not live yet.

Official binaries are distributed through GitHub Releases from `kshgr/sponsorbucks-client`. Source/dev builds are marked `build_channel=dev`; release builds are created by CI with signed/checksummed assets.

## One-command preview install

After official preview binaries are published, users can start with:

```bash
npx sponsorbucks --help
```

Or install the binary directly:

```bash
curl -fsSL https://sponsorbucks.xyz/install.sh | sh
```

```powershell
irm https://sponsorbucks.xyz/install.ps1 | iex
```

Then run:

```bash
sponsorbucks install --dry-run
sponsorbucks login
```

## Quick start

Recommended install flow:

```bash
sponsorbucks install --dry-run
```

If the dry run looks good, rerun with confirmation:

```bash
sponsorbucks install --yes
```

Then link the device:

```bash
sponsorbucks login
```

Run an agent through SponsorBucks:

```bash
sponsorbucks run --surface codex -- codex
sponsorbucks run --surface claude-code -- claude
sponsorbucks run --surface gemini-cli -- gemini
```

## What it does

- Adds local shims for detected supported tools.
- Starts a local daemon for preview integrations.
- Wraps agent sessions and emits signed heartbeat events.
- Shows a sponsored text line in terminals that support it.

## Privacy guarantee

The client never sends:

- code
- prompts
- model responses
- terminal output
- repository names
- filenames
- screenshots
- clipboard contents
- window titles

It only uses the minimum metadata needed for session management, placement delivery, and signed event integrity.

## Architecture

```text
AI agent command
   |
sponsorbucks run
   |
local client
   |-- fetch eligible text+URL ad
   |-- show sponsored line
   |-- sign event payloads
   |-- send heartbeat
   |
SponsorBucks backend
```

See:

- [ADAPTERS.md](./ADAPTERS.md)
- [SECURITY.md](./SECURITY.md)
- [PRIVACY.md](./PRIVACY.md)
- [EVENT_SCHEMA.md](./EVENT_SCHEMA.md)
- [RELEASE_SIGNING.md](./RELEASE_SIGNING.md)
- [THREAT_MODEL.md](./THREAT_MODEL.md)

