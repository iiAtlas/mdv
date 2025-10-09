#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <binary-path>" >&2
  exit 1
fi

binary_path="$1"

if [[ ! -f "$binary_path" && ! -d "$binary_path" ]]; then
  echo "notarize-macos: target not found at $binary_path" >&2
  exit 1
fi

profile_name="${MACOS_NOTARY_PROFILE_NAME:-goreleaser-notary}"

binary_dir="$(cd "$(dirname "$binary_path")" && pwd)"
binary_name="$(basename "$binary_path")"
zip_path="${binary_dir}/${binary_name}.zip"

echo "notarize-macos: creating archive ${zip_path}"

rm -f "$zip_path"
ditto -c -k --keepParent "$binary_path" "$zip_path"

echo "notarize-macos: submitting ${zip_path} for notarization (profile: ${profile_name})"

keychain_path="${KEYCHAIN_PATH:-}"
if [[ -z "$keychain_path" ]]; then
  keychain_path=$(security default-keychain -d user 2>/dev/null | sed -e 's/^ *"//' -e 's/"$//' )
fi

if [[ -n "$keychain_path" && ! -f "$keychain_path" ]]; then
  echo "notarize-macos: warning - keychain path '$keychain_path' not found, using default context" >&2
  keychain_path=""
fi

submit_args=(
  submit "$zip_path"
  --keychain-profile "$profile_name"
  --progress
  --wait
)

if [[ -n "$keychain_path" ]]; then
  submit_args+=(--keychain "$keychain_path")
fi

xcrun notarytool "${submit_args[@]}"

echo "notarize-macos: notarization complete, attempting to staple ticket"

if ! xcrun stapler staple "$binary_path"; then
  echo "notarize-macos: warning - stapler failed for $binary_name (ticket might not be stapled)" >&2
fi

rm -f "$zip_path"
echo "notarize-macos: cleaned up temporary archive"
