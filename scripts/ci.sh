#!/usr/bin/env bash
set -euo pipefail

# Docker image
GO_IMAGE="${GO_IMAGE:-golang:1.26}"

# Project root (script can be called from anywhere)
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Mode (e.g.: run/test/build/...)
MODE="${1:-run}"

# Command to run project
CMD_RUN="${CMD_RUN:-go run ./cmd/ctx-init}"

# Command to run tests
CMD_TEST="${CMD_TEST:-go test -v ./...}"

case "${MODE}" in
  run)
    CMD="${CMD_RUN}"
    ;;
  test)
    CMD="${CMD_TEST}"
    ;;
  *)
    echo "error: unsupported mode: ${MODE}" >&2
    echo "supported modes: run, test" >&2
    exit 1
    ;;
esac

echo "==> Using Docker image: ${GO_IMAGE}"
echo "==> Project root: ${ROOT_DIR}"
echo "==> Mode: ${MODE}"

docker run --rm -i \
  -v "${ROOT_DIR}:/app" \
  -w /app \
  "${GO_IMAGE}" \
  bash -c "${CMD}"
