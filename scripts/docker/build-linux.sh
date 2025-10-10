#!/usr/bin/env bash
# Helper script to build Linux GUI using Docker
# This makes it easy to test Linux builds without a Linux VM
#
# Usage:
#   ./scripts/docker/build-linux.sh              # Build both amd64 and arm64
#   ./scripts/docker/build-linux.sh amd64        # Build amd64 only
#   ./scripts/docker/build-linux.sh arm64        # Build arm64 only

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

IMAGE_NAME="mdv-linux-builder"
DOCKERFILE="scripts/docker/Dockerfile.linux-build"

# Determine which architectures to build
if [[ $# -gt 0 ]]; then
  ARCH="$1"
else
  ARCH="both"
fi

echo "=== Building Docker image for Linux GUI compilation ==="
docker build --platform linux/amd64 -f "$DOCKERFILE" -t "$IMAGE_NAME" .

build_arch() {
  local arch=$1
  local cc=$2
  local platforms="linux/${arch}"

  echo ""
  echo "=== Building Linux GUI for ${arch} in Docker container ==="
  docker run --platform linux/amd64 --rm \
    -v "$PWD:/workspace" \
    -e "WAILS_PLATFORMS=${platforms}" \
    -e "CC=${cc}" \
    "$IMAGE_NAME"
}

case "$ARCH" in
  amd64)
    build_arch "amd64" "gcc"
    ;;
  arm64)
    build_arch "arm64" "aarch64-linux-gnu-gcc"
    ;;
  both)
    build_arch "amd64" "gcc"
    build_arch "arm64" "aarch64-linux-gnu-gcc"
    ;;
  *)
    echo "Error: Unknown architecture '$ARCH'. Use 'amd64', 'arm64', or 'both'." >&2
    exit 1
    ;;
esac

echo ""
echo "=== Build complete! ==="
echo "Check cmd/mdv-gui/build/bin/ for Linux binaries"
ls -lh cmd/mdv-gui/build/bin/mdv-gui_linux_* 2>/dev/null || echo "No Linux binaries found"
