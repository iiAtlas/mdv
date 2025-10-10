#!/usr/bin/env bash
# Helper script to build Linux GUI using Docker
# This makes it easy to test Linux builds without a Linux VM

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

IMAGE_NAME="mdv-linux-builder"
DOCKERFILE="scripts/docker/Dockerfile.linux-build"

echo "=== Building Docker image for Linux GUI compilation ==="
docker build --platform linux/amd64 -f "$DOCKERFILE" -t "$IMAGE_NAME" .

echo ""
echo "=== Building Linux GUI in Docker container ==="
docker run --platform linux/amd64 --rm -v "$PWD:/workspace" "$IMAGE_NAME"

echo ""
echo "=== Build complete! ==="
echo "Check cmd/mdv-gui/build/bin/ for Linux binaries"
ls -lh cmd/mdv-gui/build/bin/mdv-gui_linux_* 2>/dev/null || echo "No Linux binaries found"
