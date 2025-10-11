#!/usr/bin/env bash
# Helper script to build Linux GUI using Docker with native builds per architecture
# Uses QEMU emulation for ARM64 builds on non-ARM hosts
#
# Usage:
#   ./scripts/docker/build-linux.sh              # Build both amd64 and arm64
#   ./scripts/docker/build-linux.sh amd64        # Build amd64 only
#   ./scripts/docker/build-linux.sh arm64        # Build arm64 only
#
# Prerequisites:
#   - Docker with buildx support (Docker Desktop >= 19.03 or docker-buildx plugin)
#   - QEMU emulation for multi-platform builds (automatically installed on Docker Desktop)

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

# Determine which architectures to build
if [[ $# -gt 0 ]]; then
  ARCH="$1"
else
  ARCH="both"
fi

# Ensure buildx is available
if ! docker buildx version >/dev/null 2>&1; then
  echo "Error: docker buildx is required but not available" >&2
  echo "Install Docker Desktop or the docker-buildx plugin" >&2
  exit 1
fi

# Create buildx builder if it doesn't exist
if ! docker buildx inspect mdv-builder >/dev/null 2>&1; then
  echo "=== Creating multi-platform buildx builder ==="
  docker buildx create --name mdv-builder --use --bootstrap
else
  docker buildx use mdv-builder
fi

build_arch() {
  local arch=$1
  local platform="linux/${arch}"
  local dockerfile="scripts/docker/Dockerfile.linux-${arch}"
  local image_name="mdv-linux-${arch}"

  echo ""
  echo "=== Building Docker image for ${platform} ==="
  docker buildx build \
    --platform "${platform}" \
    --load \
    -f "${dockerfile}" \
    -t "${image_name}" \
    .

  echo ""
  echo "=== Building Linux GUI for ${arch} in Docker container ==="
  docker run \
    --platform "${platform}" \
    --rm \
    -v "$PWD:/workspace" \
    "${image_name}"
}

case "$ARCH" in
  amd64)
    build_arch "amd64"
    ;;
  arm64)
    echo "NOTE: ARM64 build uses QEMU emulation and will be slower than native builds"
    build_arch "arm64"
    ;;
  both)
    build_arch "amd64"
    echo ""
    echo "NOTE: ARM64 build uses QEMU emulation and will be slower than native builds"
    build_arch "arm64"
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
