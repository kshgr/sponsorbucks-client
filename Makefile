SHELL := /bin/sh

VERSION ?= 1.0.0-preview
BUILD_ID ?= local
BUILD_CHANNEL ?= dev

.PHONY: build test clean release

build:
	go build ./cmd/sponsorbucks

test:
	go test ./...

release:
	rm -rf dist
	mkdir -p dist
	GOOS=windows GOARCH=amd64 VERSION=$(VERSION) BUILD_ID=$(BUILD_ID) BUILD_CHANNEL=release sh ./scripts/build-release.sh
	GOOS=linux GOARCH=amd64 VERSION=$(VERSION) BUILD_ID=$(BUILD_ID) BUILD_CHANNEL=release sh ./scripts/build-release.sh
	GOOS=darwin GOARCH=amd64 VERSION=$(VERSION) BUILD_ID=$(BUILD_ID) BUILD_CHANNEL=release sh ./scripts/build-release.sh
	GOOS=darwin GOARCH=arm64 VERSION=$(VERSION) BUILD_ID=$(BUILD_ID) BUILD_CHANNEL=release sh ./scripts/build-release.sh
	cd dist && sha256sum sponsorbucks-* > checksums.txt

clean:
	rm -rf dist sponsorbucks sponsorbucks.exe
