#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "usage: $0 <path-to-binary-or-bundle> [entitlements]" >&2
  exit 1
fi

binary_path="$1"
entitlements="${2:-${MACOS_SIGN_ENTITLEMENTS:-}}"

if [[ ! -f "$binary_path" && ! -d "$binary_path" ]]; then
  echo "sign-macos: target not found at $binary_path" >&2
  exit 1
fi

codesign_identity="${CODESIGN_IDENTITY:-${MACOS_SIGN_IDENTITY:-}}"
if [[ -z "$codesign_identity" ]]; then
  echo "sign-macos: CODESIGN_IDENTITY env var is required" >&2
  exit 1
fi

if [[ -n "$entitlements" ]]; then
  if [[ ! -f "$entitlements" ]]; then
    echo "sign-macos: entitlements file not found at $entitlements" >&2
    exit 1
  fi
fi

echo "sign-macos: signing $(basename "$binary_path") with identity '$codesign_identity'"

codesign_args=(
  --force
  --options runtime
  --timestamp
  --sign "$codesign_identity"
)

if [[ -n "$entitlements" ]]; then
  codesign_args+=(--entitlements "$entitlements")
fi

if [[ -d "$binary_path" ]]; then
  codesign_args+=(--deep)
fi

codesign "${codesign_args[@]}" "$binary_path"

codesign --verify --deep --strict "$binary_path"

echo "sign-macos: codesign complete"
