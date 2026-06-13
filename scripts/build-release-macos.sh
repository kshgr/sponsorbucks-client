#!/bin/sh
set -eu

GOOS=darwin GOARCH="${GOARCH:-amd64}" OUT_DIR="${OUT_DIR:-dist}" sh ./scripts/build-release.sh
