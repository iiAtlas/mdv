#!/usr/bin/env bash

# Downloads the latest GitHub release artifacts and verifies macOS signing.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

if ! command -v gh >/dev/null 2>&1; then
  echo "verify-latest-release: gh CLI is required" >&2
  exit 1
fi

if ! command -v codesign >/dev/null 2>&1 || ! command -v spctl >/dev/null 2>&1; then
  echo "verify-latest-release: macOS codesign/spctl tools are required" >&2
  exit 1
fi

REPO_NAME="${1:-}"
if [[ -z "$REPO_NAME" ]]; then
  REPO_NAME="$(gh repo view --json nameWithOwner --jq .nameWithOwner 2>/dev/null || true)"
fi
if [[ -z "$REPO_NAME" ]]; then
  origin_url="$(git -C "$REPO_ROOT" config --get remote.origin.url || true)"
  if [[ "$origin_url" =~ github.com[:/](.+/.+)(\.git)?$ ]]; then
    REPO_NAME="${BASH_REMATCH[1]}"
  fi
fi
if [[ -z "$REPO_NAME" ]]; then
  echo "verify-latest-release: unable to determine repository name, pass it explicitly (e.g. ./verify-latest-release.sh owner/repo)" >&2
  exit 1
fi

LATEST_TAG="$(gh release view --repo "$REPO_NAME" --json tagName --jq .tagName)"
echo "verify-latest-release: latest tag is ${LATEST_TAG}"

WORK_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

echo "verify-latest-release: downloading macOS artifacts to ${WORK_DIR}"
gh release download "$LATEST_TAG" \
  --repo "$REPO_NAME" \
  --dir "$WORK_DIR" \
  --pattern "mdv_*.tar.gz" \
  --pattern "mdv-.gui_*.tar.gz"

check_codesign() {
  local target="$1"
  local label="$2"

  if [[ ! -e "$target" ]]; then
    echo "verify: missing ${label} (${target})"
    return 1
  fi

  echo "verify: codesign --verify --deep --strict ${label}"
  codesign --verify --deep --strict "$target"

  echo "verify: spctl --assess ${label}"
  set +e
  spctl --assess --type execute --verbose "$target"
  local status=$?
  set -e

  if [[ $status -eq 3 ]]; then
    echo "verify: spctl informational status 3 for ${label} (expected for non-app binaries)"
  elif [[ $status -ne 0 ]]; then
    echo "verify: spctl returned ${status} for ${label}"
    return "$status"
  fi
}

status=0

for tarball in "$WORK_DIR"/mdv_*_darwin_*.tar.gz; do
  [[ -f "$tarball" ]] || continue
  basename="$(basename "$tarball")"
  if [[ "$basename" == mdv-.gui_* ]]; then
    continue
  fi
  extract_dir="$(mktemp -d "${WORK_DIR}/cli-XXXX")"
  tar -xzf "$tarball" -C "$extract_dir"
  cli_bin="${extract_dir}/mdv"
  if ! check_codesign "$cli_bin" "CLI $basename"; then
    status=1
  fi
done

for tarball in "$WORK_DIR"/mdv-.gui_*_darwin_*.tar.gz; do
  [[ -f "$tarball" ]] || continue
  basename="$(basename "$tarball")"
  extract_dir="$(mktemp -d "${WORK_DIR}/gui-XXXX")"
  tar -xzf "$tarball" -C "$extract_dir"
  gui_bin="${extract_dir}/mdv-gui"
  gui_app="${extract_dir}/mdv-gui.app"
  if ! check_codesign "$gui_bin" "GUI binary $basename"; then
    status=1
  fi
  if ! check_codesign "$gui_app" "GUI app $basename"; then
    status=1
  fi
  if [[ -d "$gui_app" ]]; then
    inner_bin="${gui_app}/Contents/MacOS/mdv-gui"
    if ! check_codesign "$inner_bin" "GUI app executable $basename"; then
      status=1
    fi
  fi
done

if [[ $status -ne 0 ]]; then
  echo "verify-latest-release: verification failed"
  exit "$status"
fi

echo "verify-latest-release: all macOS artifacts verified for ${LATEST_TAG}"
