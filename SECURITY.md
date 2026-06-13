# Security

SponsorBucks client v0.3 is designed around a narrow local-trust boundary.

## Non-negotiables

- No code collection.
- No prompt collection.
- No response collection.
- No terminal output collection.
- No repo-name collection.
- No filename collection.
- No screenshot collection.
- No clipboard collection.
- No window-title collection.
- No env-var collection.
- No secret collection.
- No remote executable code.
- No auto-update in v0.

## Release expectations

- Official binaries must be signed and checksummed.
- Dev/source builds must be labeled clearly in build metadata.
- The backend may treat dev builds as non-payable or limited.

## Local attack surface

- `sponsorbucks daemon` listens only on localhost.
- Shell shims are text files that forward to the local CLI.
- `sponsorbucks run` should not log command arguments after `--` by default.
