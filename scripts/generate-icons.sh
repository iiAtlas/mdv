#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 path/to/source-icon.png" >&2
  exit 1
fi

SOURCE_ICON="$1"
if [[ ! -f "$SOURCE_ICON" ]]; then
  echo "Error: file not found: $SOURCE_ICON" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ASSET_DIR="${REPO_ROOT}/cmd/mdv-gui/frontend/assets"
BUILD_DIR="${REPO_ROOT}/cmd/mdv-gui/build"
WINDOWS_DIR="${BUILD_DIR}/windows"

mkdir -p "${ASSET_DIR}" "${WINDOWS_DIR}"

copy_icon() {
  local dest="$1"
  if [[ -e "${dest}" ]] && [[ "${SOURCE_ICON}" -ef "${dest}" ]]; then
    return
  fi
  install -m 0644 "${SOURCE_ICON}" "${dest}"
}

resize_image() {
  local size="$1"
  local dest="$2"

  if command -v sips >/dev/null 2>&1; then
    sips -s format png -z "${size}" "${size}" "${SOURCE_ICON}" --out "${dest}" >/dev/null
  elif command -v ffmpeg >/dev/null 2>&1; then
    ffmpeg -loglevel error -y -i "${SOURCE_ICON}" -vf "scale=${size}:${size}:flags=lanczos" -frames:v 1 "${dest}"
  else
    echo "Error: neither 'sips' nor 'ffmpeg' is available to resize images." >&2
    exit 1
  fi
}

TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t mdvicon)"
cleanup() {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT
export TMP_DIR
export REPO_ROOT

# Update canonical PNG copies used by Wails and the repo.
copy_icon "${ASSET_DIR}/mdv-app-icon.png"
copy_icon "${BUILD_DIR}/appicon.png"

# Prepare PNG variants for the Windows .ico bundle.
SIZES=(16 24 32 48 64 128 256)
for size in "${SIZES[@]}"; do
  resize_image "${size}" "${TMP_DIR}/${size}.png"
done

# Generate ICO file from the prepared PNGs.
python3 - <<'PY'
import os
import struct

tmp_dir = os.environ["TMP_DIR"]
repo_root = os.environ["REPO_ROOT"]

sizes = [16, 24, 32, 48, 64, 128, 256]
entries = []
payloads = []
offset = 6 + len(sizes) * 16
for size in sizes:
    path = os.path.join(tmp_dir, f"{size}.png")
    with open(path, "rb") as f:
        data = f.read()
    width = size if size < 256 else 0
    height = size if size < 256 else 0
    entries.append(struct.pack("<BBBBHHII", width, height, 0, 0, 1, 32, len(data), offset))
    payloads.append(data)
    offset += len(data)

header = struct.pack("<HHH", 0, 1, len(sizes))
blob = header + b"".join(entries) + b"".join(payloads)

targets = [
    os.path.join(repo_root, "cmd/mdv-gui/frontend/assets/mdv-app-icon.ico"),
    os.path.join(repo_root, "cmd/mdv-gui/build/windows/icon.ico"),
]
for target in targets:
    os.makedirs(os.path.dirname(target), exist_ok=True)
    with open(target, "wb") as f:
        f.write(blob)
PY

echo "Updated icons:"
echo "  - ${ASSET_DIR}/mdv-app-icon.png"
echo "  - ${ASSET_DIR}/mdv-app-icon.ico"
echo "  - ${BUILD_DIR}/appicon.png"
echo "  - ${WINDOWS_DIR}/icon.ico"
