# Agent wrapper examples

## Codex

```bash
sponsorbucks run --surface codex -- codex
```

## Claude Code

```bash
sponsorbucks run --surface claude-code -- claude
```

## Pi

```bash
sponsorbucks run --surface pi -- pi
```

## OpenClaw

```bash
sponsorbucks run --surface openclaw -- openclaw
```

## Generic command

```bash
sponsorbucks run --surface terminal -- your-agent-command --with --args
```

## Demo mode without backend

```bash
sponsorbucks run --demo --surface codex -- ping google.com
```

Windows PowerShell example:

```powershell
sponsorbucks run --demo --surface codex -- ping 127.0.0.1
```
