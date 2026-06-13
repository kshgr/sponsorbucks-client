# SponsorBucks Client Architecture

## Components

### Go CLI

Binary name: `sponsorbucks`

Commands:

- `login` - opens browser for device linking.
- `status` - shows linked user/device state.
- `logout` - clears local session.
- `logs` - prints local SponsorBucks logs only.
- `install` / `uninstall` - manage local shims and shell PATH integration.
- `pause` / `resume` - globally pause or resume placements.
- `enable <tool>` / `disable <tool>` - toggle a specific tool.
- `run --surface <surface> -- <command>` - wraps an agent command.
- `config set-api <url>` - sets the backend base URL.
- `config show` - prints redacted local config.
- `debug event [--json|--explain]` - prints a sample event payload without sending.
- `privacy` - prints the privacy summary.
- `doctor` - checks local readiness.

### Local config

Stored at:

- macOS/Linux: `~/.sponsorbucks/config.json`
- Windows: `%USERPROFILE%\.sponsorbucks\config.json`

Contains:

- API base URL
- device ID
- device private key
- linked user token/session token
- client preferences

Never store advertiser secrets, service-role keys, or admin tokens locally.

### Device identity

On first launch:

- Generate an Ed25519 key pair.
- Store the private key locally.
- Send the public key during device linking.
- Sign every event body with the private key.
- Backend verifies the signature using the registered public key.

### Run wrapper

For v0:

- Start the child process.
- Mark `session_start`.
- Every 5 seconds:
  - fetch or rotate an eligible ad;
  - show one sponsored line;
  - submit heartbeat.
- On process exit:
  - submit `session_end`.

### Attention signals

Version one should keep this intentionally conservative:

- screen unlocked: best-effort OS-specific check.
- recent input: best-effort OS-specific check; initially bucketed and allowed to be unknown.
- foreground supported surface: true when wrapper controls the terminal or a preview adapter displays status.
- one paid placement at a time: enforced by backend session state and local mutex.

Do not implement screenshot-based verification.

### Adapters

Preview adapters are display-only:

- shell shims for detected tools;
- local daemon for localhost integrations;
- no workspace reads;
- no editor content capture;
- no terminal-output capture.

