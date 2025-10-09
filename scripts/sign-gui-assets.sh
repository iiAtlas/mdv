#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

codesign_identity="${CODESIGN_IDENTITY:-${MACOS_SIGN_IDENTITY:-}}"
if [[ -z "$codesign_identity" ]]; then
  echo "sign-gui-assets: CODESIGN_IDENTITY or MACOS_SIGN_IDENTITY must be set" >&2
  exit 1
fi

gui_entitlements="${MACOS_GUI_ENTITLEMENTS:-${repo_root}/cmd/mdv-gui/entitlements.plist}"
if [[ -n "$gui_entitlements" && ! -f "$gui_entitlements" ]]; then
  echo "sign-gui-assets: GUI entitlements file not found at $gui_entitlements" >&2
  exit 1
fi

artifact_root="${repo_root}/cmd/mdv-gui/build/bin"
shopt -s nullglob
targets=( "${artifact_root}"/mdv-gui_* )
shopt -u nullglob

if [[ ${#targets[@]} -eq 0 ]]; then
  echo "sign-gui-assets: no GUI artifacts found, skipping" >&2
  exit 0
fi

for target_dir in "${targets[@]}"; do
  [[ -d "$target_dir" ]] || continue

  arch_label="${target_dir#${repo_root}/}"
  bin_path="${target_dir}/mdv-gui"
  app_path="${target_dir}/mdv-gui.app"

  if [[ -f "$bin_path" ]]; then
    echo "sign-gui-assets: signing binary ${bin_path#${repo_root}/}"
    "${script_dir}/sign-macos.sh" "$bin_path" "$gui_entitlements"
    echo "sign-gui-assets: notarizing binary ${bin_path#${repo_root}/}"
    "${script_dir}/notarize-macos.sh" "$bin_path"
  fi

  if [[ -d "$app_path" ]]; then
    inner_binary="${app_path}/Contents/MacOS/mdv-gui"
    if [[ -f "$inner_binary" ]]; then
      echo "sign-gui-assets: signing bundle executable ${inner_binary#${repo_root}/}"
      "${script_dir}/sign-macos.sh" "$inner_binary" "$gui_entitlements"
    fi
    echo "sign-gui-assets: signing app bundle ${app_path#${repo_root}/}"
    "${script_dir}/sign-macos.sh" "$app_path" "$gui_entitlements"
    echo "sign-gui-assets: notarizing app bundle ${app_path#${repo_root}/}"
    "${script_dir}/notarize-macos.sh" "$app_path"
  fi

  echo "sign-gui-assets: finished ${arch_label}"
done

echo "sign-gui-assets: completed"
