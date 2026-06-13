# Release Signing

Official SponsorBucks releases must ship with:

- a signed release artifact
- a checksum manifest
- a published client version
- a published build ID
- a published build channel

Recommended release flow:

1. Build the binary with explicit build metadata.
2. Generate checksums for every release artifact.
3. Sign the artifact or checksum manifest with the release key.
4. Publish the signature alongside the release.

Source builds should default to `build_channel=dev` unless flags override it.

