#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <binary-path>" >&2
  exit 1
fi

binary_path="$1"

if [[ ! -f "$binary_path" ]]; then
  echo "sign-macos: binary not found at $binary_path" >&2
  exit 1
fi

codesign_identity="${CODESIGN_IDENTITY:-${MACOS_SIGN_IDENTITY:-}}"
if [[ -z "$codesign_identity" ]]; then
  echo "sign-macos: CODESIGN_IDENTITY env var is required" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"
entitlements="${repo_root}/cmd/mdv-gui/entitlements.plist"

if [[ ! -f "$entitlements" ]]; then
  echo "sign-macos: entitlements file not found at $entitlements" >&2
  exit 1
fi

echo "sign-macos: signing $(basename "$binary_path") with identity '$codesign_identity'"

codesign --force \
  --options runtime \
  --timestamp \
  --entitlements "$entitlements" \
  --sign "$codesign_identity" \
  "$binary_path"

codesign --verify --deep --strict "$binary_path"

echo "sign-macos: codesign complete"
