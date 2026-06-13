# Security and Privacy Model

## Non-negotiables

The client never reads or sends:

- source code
- prompts
- model responses
- terminal output
- filenames
- repository names
- screenshots
- clipboard contents
- window titles

## Threat model

### Fraud

Risks:

- automated fake sessions
- replayed heartbeats
- many accounts per device
- referral rings
- hidden/minimized terminals
- unattended overnight agents
- self-clicks

Controls:

- device key pair and signed events
- sequence numbers
- one earning placement per user at a time
- visible duration threshold
- daily caps
- delayed earnings clearance
- referral review
- manual admin review
- remote client version allowlist/kill switch

### Privacy

Risks:

- users fear code inspection
- enterprise developers cannot install invasive tools

Controls:

- open-source client
- published event schema
- no terminal capture
- no screenshots
- no workspace file reads
- no file-system scanning
- no clipboard access
- minimal local logs
- clear uninstall path

## Local logs

Local logs should include:

- timestamps
- endpoint status
- session IDs
- surface name
- error messages

Local logs should not include:

- command arguments after `--` by default
- terminal output
- current directory
- repo name

If debug mode records more detail, it must be explicit and opt-in.
