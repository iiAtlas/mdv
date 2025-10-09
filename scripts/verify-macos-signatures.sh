#!/usr/bin/env bash

# Verifies codesign/spctl status for the locally built macOS artifacts.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

check_codesign() {
  local target="$1"
  local label="$2"

  if [[ ! -e "$target" ]]; then
    echo "verify: skipping ${label} (missing: ${target})"
    return
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
  fi
}

echo "verify: checking macOS CLI binaries in dist/"
while IFS= read -r -d '' path; do
  label="CLI $(basename "$(dirname "$path")")"
  check_codesign "$path" "$label"
done < <(find "${ROOT_DIR}/dist" -maxdepth 2 -type f -name mdv -path "*darwin*" -print0 2>/dev/null)

echo "verify: checking GUI binaries and app bundles"
for artifact_dir in "${ROOT_DIR}"/cmd/mdv-gui/build/bin/mdv-gui_*; do
  [[ -d "$artifact_dir" ]] || continue

  gui_bin="${artifact_dir}/mdv-gui"
  gui_app="${artifact_dir}/mdv-gui.app"

  check_codesign "$gui_bin" "GUI binary ${artifact_dir##*/}"
  check_codesign "$gui_app" "GUI app ${artifact_dir##*/}"

  if [[ -d "$gui_app" ]]; then
    inner_bin="${gui_app}/Contents/MacOS/mdv-gui"
    check_codesign "$inner_bin" "GUI app executable ${artifact_dir##*/}"
  fi

done

echo "verify: (optional) checking extracted tarballs"
for tarball in "${ROOT_DIR}"/dist/mdv_*_darwin_*.tar.gz; do
  [[ -f "$tarball" ]] || continue
  tmpdir="$(mktemp -d)"
  tar -xzf "$tarball" -C "$tmpdir"
  extracted="${tmpdir}/mdv"
  if [[ -f "$extracted" ]]; then
    check_codesign "$extracted" "tarball $(basename "$tarball")"
  fi
  rm -rf "$tmpdir"
done

echo "verify: completed"
