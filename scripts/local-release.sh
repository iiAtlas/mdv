#!/usr/bin/env bash

# Convenience wrapper to run a local GoReleaser snapshot with the signing
# environment loaded from .env.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "local-release: missing .env file at ${ENV_FILE}" >&2
  exit 1
fi

if ! command -v goreleaser >/dev/null 2>&1; then
  echo "local-release: goreleaser is not installed or not in PATH" >&2
  exit 1
fi

pushd "$ROOT_DIR" >/dev/null

echo "local-release: sourcing .env and running goreleaser snapshot"
set -a
source "$ENV_FILE"
set +a

goreleaser release --snapshot --skip=publish --clean "$@"

popd >/dev/null
