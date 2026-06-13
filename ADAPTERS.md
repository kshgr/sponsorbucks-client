# Adapters

SponsorBucks v0.3 supports these local adapters and compatibility surfaces:

| Surface | Tool command | Status |
| --- | --- | --- |
| Codex | `codex` | tested |
| Claude Code | `claude` | preview |
| Pi | `pi` | preview |
| Aider | `aider` | planned |
| OpenCode | `opencode` | planned |
| Gemini CLI | `gemini` | preview |
| Generic terminal | any wrapped command | tested |

Compatibility notes:

- `claude-code` maps to the `claude` command for backward compatibility.
- `gemini-cli` maps to the `gemini` command for backward compatibility.
- adapters must not read workspace files.
- adapters must not collect terminal output.
- adapters must not patch third-party tool files.
- adapters must forward the real command unchanged after `--`.

