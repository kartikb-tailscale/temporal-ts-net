#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
PORT_BASE="$((10000 + RANDOM % 1000))"
GRPC_PORT="${GRPC_PORT:-${PORT_BASE}}"
UI_PORT="${UI_PORT:-$((GRPC_PORT + 1000))}"
METRICS_PORT="${METRICS_PORT:-$((GRPC_PORT + 2000))}"
ADDRESS="127.0.0.1:${GRPC_PORT}"
CODEC_PORT="${CODEC_PORT:-8081}"
CODEC_ENDPOINT="http://127.0.0.1:${CODEC_PORT}"

cleanup() {
  jobs -pr | xargs -I{} kill {} 2>/dev/null || true
}
trap cleanup EXIT INT TERM

cd "${ROOT_DIR}"

echo "Starting Temporal dev server through extension (with built-in codec server)"
temporal start-dev \
  --ip 127.0.0.1 \
  --port "${GRPC_PORT}" \
  --codec-port "${CODEC_PORT}" \
  --ui-port "${UI_PORT}" \
  --metrics-port "${METRICS_PORT}" \
  > /tmp/temporal-start-dev-demo.log 2>&1 &
echo "Using codec endpoint: ${CODEC_ENDPOINT}"

for _ in $(seq 1 60); do
  if temporal workflow list --address "${ADDRESS}" --limit 1 >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done

echo "Starting SDK worker"
go run ./worker \
  --address "${ADDRESS}" \
  --codec-endpoint "${CODEC_ENDPOINT}" > /tmp/worker-demo.log 2>&1 &
sleep 2

echo "Starting workflow"
go run ./starter \
  --address "${ADDRESS}" \
  --name "codec-demo" \
  --codec-endpoint "${CODEC_ENDPOINT}"

echo "Demo complete"
echo "Temporal logs: /tmp/temporal-start-dev-demo.log"
echo "Worker logs:   /tmp/worker-demo.log"
