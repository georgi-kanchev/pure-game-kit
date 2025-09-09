#!/usr/bin/env bash
set -euo pipefail

echo "Building Linux app..."
go build -ldflags="-s -w" -o app .
echo "Build complete: ./app"
