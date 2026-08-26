#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_PID=""
FRONTEND_PID=""

cleanup() {
	echo ""
	echo "Deteniendo servicios..."
	[[ -n "$FRONTEND_PID" ]] && kill "$FRONTEND_PID" 2>/dev/null || true
	[[ -n "$BACKEND_PID" ]] && kill "$BACKEND_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

echo "▶ Levantando backend (Go) en :8080"
(cd "$ROOT_DIR/backend" && go run .) &
BACKEND_PID=$!

echo "▶ Levantando frontend (SvelteKit)"
(cd "$ROOT_DIR/frontend" && pnpm dev) &
FRONTEND_PID=$!

wait "$BACKEND_PID" "$FRONTEND_PID"