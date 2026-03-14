#!/usr/bin/env bash
# Install go-gen-r for macOS / Linux: build and copy to /usr/local/bin.
# For Windows, use install.ps1 instead.

set -e
cd "$(dirname "$0")"
go build -o go-gen-r ./cmd/go-gen-r
mv -f go-gen-r /usr/local/bin/go-gen-r
echo "Installed: /usr/local/bin/go-gen-r"