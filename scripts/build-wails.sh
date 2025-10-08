#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Wails can only produce native bundles reliably for macOS from macOS;
# cross-compiling GUI binaries for other OSes is handled on their respective
# builders. Keep the defaults focused on the supported macOS targets.
DEFAULT_PLATFORMS="darwin/amd64,darwin/arm64"
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

# Use local cache if possible (for local builds), fall back to default (for CI)
if mkdir -p "${ROOT_DIR}/.cache/go-build" "${ROOT_DIR}/.cache/go-mod" 2>/dev/null; then
  export GOCACHE="${ROOT_DIR}/.cache/go-build"
  export GOMODCACHE="${ROOT_DIR}/.cache/go-mod"
fi

pushd "$ROOT_DIR/cmd/mdv-gui" >/dev/null

echo "[wails] building GUI for platforms: ${PLATFORMS}" >&2
wails build -clean -platform "${PLATFORMS}" -ldflags "${LDFLAGS}" "$@"

# Normalise the bundle locations so downstream tooling (GoReleaser) can pick
# them up consistently.
BIN_DIR="$ROOT_DIR/cmd/mdv-gui/build/bin"
for platform in ${PLATFORMS//,/ }; do
  IFS='/' read -r os arch <<<"$platform"
  case "$os" in
    darwin)
      bundle="$BIN_DIR/mdv-gui-${arch}.app"
      if [[ ! -d "$bundle" ]]; then
        # Wails emits unsuffixed bundle names when targeting a single arch.
        bundle="$BIN_DIR/mdv-gui.app"
      fi
      if [[ -d "$bundle" ]]; then
        dest="$BIN_DIR/mdv-gui_${os}_${arch}"
        rm -rf "$dest"
        mkdir -p "$dest"
        cp -R "$bundle" "$dest/mdv-gui.app"
        cp "$bundle/Contents/MacOS/mdv-gui" "$dest/mdv-gui"
      else
        echo "[wails] warning: expected bundle ${bundle} not found" >&2
      fi
      ;;
    windows)
      dest="$BIN_DIR/mdv-gui_${os}_${arch}"
      rm -rf "$dest"
      mkdir -p "$dest"
      for candidate in \
        "$BIN_DIR/mdv-gui-${arch}.exe" \
        "$BIN_DIR/mdv-gui-${os}-${arch}.exe" \
        "$BIN_DIR/mdv-gui.exe"; do
        if [[ -f "$candidate" ]]; then
          cp "$candidate" "$dest/mdv-gui.exe"
          break
        fi
      done
      if [[ ! -f "$dest/mdv-gui.exe" ]]; then
        echo "[wails] warning: expected binary for ${os}/${arch} not found" >&2
        rmdir "$dest" 2>/dev/null || true
      fi
      ;;
    linux)
      dest="$BIN_DIR/mdv-gui_${os}_${arch}"
      rm -rf "$dest"
      mkdir -p "$dest"
      for candidate in \
        "$BIN_DIR/mdv-gui-${os}-${arch}" \
        "$BIN_DIR/mdv-gui-${arch}" \
        "$BIN_DIR/mdv-gui"; do
        if [[ -f "$candidate" ]]; then
          cp "$candidate" "$dest/mdv-gui"
          break
        fi
      done
      if [[ ! -f "$dest/mdv-gui" ]]; then
        echo "[wails] warning: expected binary for ${os}/${arch} not found" >&2
        rmdir "$dest" 2>/dev/null || true
      fi
      ;;
    *)
      echo "[wails] skipping bundle normalisation for ${os}/${arch}" >&2
      ;;
  esac
done

popd >/dev/null
