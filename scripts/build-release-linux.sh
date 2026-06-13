#!/bin/sh
set -eu

GOOS=linux GOARCH="${GOARCH:-amd64}" OUT_DIR="${OUT_DIR:-dist}" sh ./scripts/build-release.sh
