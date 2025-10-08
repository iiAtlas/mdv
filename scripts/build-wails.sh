#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DEFAULT_PLATFORMS="darwin/amd64,darwin/arm64,linux/amd64,linux/arm64,windows/amd64"
if [[ $# -gt 0 ]]; then
  PLATFORMS=$1
  shift
else
  PLATFORMS=$DEFAULT_PLATFORMS
fi

VERSION=${GORELEASER_CURRENT_TAG:-${VERSION:-dev}}
COMMIT=${GORELEASER_COMMIT:-$(cd "$ROOT_DIR" && git rev-parse HEAD)}
DATE=${GORELEASER_DATE:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

export GOCACHE="${ROOT_DIR}/.cache/go-build"
export GOMODCACHE="${ROOT_DIR}/.cache/go-mod"
mkdir -p "$GOCACHE" "$GOMODCACHE"

pushd "$ROOT_DIR/cmd/mdv-gui" >/dev/null

echo "[wails] building GUI for platforms: ${PLATFORMS}" >&2
wails build -clean -platform "${PLATFORMS}" -ldflags "${LDFLAGS}" "$@"

popd >/dev/null
