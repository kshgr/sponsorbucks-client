#!/bin/sh
set -eu

GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
OUT_DIR="${OUT_DIR:-dist}"
VERSION="${VERSION:-1.0.0-preview}"
BUILD_ID="${BUILD_ID:-manual}"
BUILD_CHANNEL="${BUILD_CHANNEL:-release}"

mkdir -p "$OUT_DIR"
ext=""
if [ "$GOOS" = "windows" ]; then ext=".exe"; fi
out="$OUT_DIR/sponsorbucks-${GOOS}-${GOARCH}${ext}"
ldflags="-s -w -X sponsorbucks-client/internal/buildinfo.Version=$VERSION -X sponsorbucks-client/internal/buildinfo.BuildID=$BUILD_ID -X sponsorbucks-client/internal/buildinfo.BuildChannel=$BUILD_CHANNEL"
GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "$ldflags" -o "$out" ./cmd/sponsorbucks
