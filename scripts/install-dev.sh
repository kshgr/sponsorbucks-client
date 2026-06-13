#!/usr/bin/env bash
set -euo pipefail

go build -o sponsorbucks ./cmd/sponsorbucks
mkdir -p "$HOME/.local/bin"
mv sponsorbucks "$HOME/.local/bin/sponsorbucks"

echo "Installed sponsorbucks to $HOME/.local/bin/sponsorbucks"
echo "Make sure $HOME/.local/bin is in your PATH."
