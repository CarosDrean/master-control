#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_PID=""
FRONTEND_PID=""

LAN_IP="$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || true)"

cleanup() {
	echo ""
	echo "Deteniendo servicios..."
	[[ -n "$FRONTEND_PID" ]] && kill "$FRONTEND_PID" 2>/dev/null || true
	[[ -n "$BACKEND_PID" ]] && kill "$BACKEND_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

if [[ -n "$LAN_IP" ]]; then
	echo ""
	echo "════════════════════════════════════════════════════"
	echo "  Desde tu tablet/celular conectate a:"
	echo "    UI    → http://${LAN_IP}:5173"
	echo "    API   → http://${LAN_IP}:8080/master-control"
	echo "  (Tu Mac y el dispositivo deben estar en la misma red)"
	echo "════════════════════════════════════════════════════"
	echo ""
fi

echo "▶ Levantando backend (Go) en :8080"
(cd "$ROOT_DIR/backend" && go run .) &
BACKEND_PID=$!

echo "▶ Levantando frontend (SvelteKit) en todas las interfaces"
(cd "$ROOT_DIR/frontend" && pnpm dev --host) &
FRONTEND_PID=$!

wait "$BACKEND_PID" "$FRONTEND_PID"