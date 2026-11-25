#!/usr/bin/env bash
# Levanta backend + frontend.
#
#   ./dev.sh               solo esta Mac (127.0.0.1) — por defecto
#   ./dev.sh --tailscale   tu tailnet (iPad/iPhone con Tailscale, tráfico cifrado) — recomendado
#   ./dev.sh --lan         la red local (HTTP sin cifrar: el token viaja en claro por el WiFi)
#
# El token (MC_TOKEN) se genera una vez en backend/.env y se reutiliza. La URL que se imprime
# lo lleva en el fragmento (#token=…): ábrela una vez en cada dispositivo y queda recordado.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_PID=""
FRONTEND_PID=""
ENV_FILE="$ROOT_DIR/backend/.env"
UI_PORT=5173

MODE="local"
case "${1:-}" in
	"" | --local) MODE="local" ;;
	--tailscale) MODE="tailscale" ;;
	--lan) MODE="lan" ;;
	-h | --help)
		sed -n '2,9p' "$0" | sed 's/^# \{0,1\}//'
		exit 0
		;;
	*)
		echo "Opción desconocida: $1 (usa --local, --tailscale o --lan)" >&2
		exit 2
		;;
esac

if [[ ! -f "$ENV_FILE" ]]; then
	echo "▶ Creando $ENV_FILE con valores por defecto"
	cat > "$ENV_FILE" <<'EOF'
GIN_MODE=release
APP_NAME=master-control
BUCKET_NAME=actions
EOF
fi

# Token persistente: se genera una sola vez y nunca se sobrescribe.
if ! grep -qE '^MC_TOKEN=.+' "$ENV_FILE"; then
	echo "▶ Generando MC_TOKEN en $ENV_FILE"
	printf 'MC_TOKEN=%s\n' "$(openssl rand -hex 24)" >> "$ENV_FILE"
fi
chmod 600 "$ENV_FILE"
TOKEN="$(grep -E '^MC_TOKEN=' "$ENV_FILE" | tail -1 | cut -d= -f2-)"

tailscale_ip() {
	local bin
	for bin in tailscale /Applications/Tailscale.app/Contents/MacOS/Tailscale; do
		if command -v "$bin" >/dev/null 2>&1; then
			"$bin" ip -4 2>/dev/null | head -1
			return 0
		fi
	done
	return 1
}

case "$MODE" in
	local)
		BIND_HOST="127.0.0.1"
		PUBLIC_HOST="localhost"
		;;
	tailscale)
		BIND_HOST="$(tailscale_ip || true)"
		if [[ -z "$BIND_HOST" ]]; then
			echo "✕ Tailscale no está activo (no hay IP de tailnet). Ábrelo e inicia sesión." >&2
			exit 1
		fi
		PUBLIC_HOST="$BIND_HOST"
		;;
	lan)
		BIND_HOST="0.0.0.0"
		PUBLIC_HOST="$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null || true)"
		if [[ -z "$PUBLIC_HOST" ]]; then
			echo "✕ No se detectó una IP de red local." >&2
			exit 1
		fi
		;;
esac

UI_URL="http://${PUBLIC_HOST}:${UI_PORT}"
# El backend solo acepta peticiones de la UI servida en este host (más localhost para desarrollo).
ORIGINS="${UI_URL},http://localhost:${UI_PORT},http://127.0.0.1:${UI_PORT}"

cleanup() {
	echo ""
	echo "Deteniendo servicios..."
	[[ -n "$FRONTEND_PID" ]] && kill "$FRONTEND_PID" 2>/dev/null || true
	[[ -n "$BACKEND_PID" ]] && kill "$BACKEND_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

LOGIN_URL="${UI_URL}/#token=${TOKEN}"

echo ""
echo "════════════════════════════════════════════════════"
echo "  Modo: ${MODE}"
echo "  Abre (una vez por dispositivo):"
echo "    ${LOGIN_URL}"
if [[ "$MODE" == "lan" ]]; then
	echo ""
	echo "  ⚠ HTTP sin cifrar en la red local: cualquiera en este WiFi"
	echo "    puede ver el token. Prefiere ./dev.sh --tailscale."
fi
echo "════════════════════════════════════════════════════"
if [[ "$MODE" != "local" ]]; then
	if command -v qrencode >/dev/null 2>&1; then
		qrencode -t ansiutf8 "$LOGIN_URL"
	else
		echo "  (instala qrencode para ver un QR: brew install qrencode)"
	fi
fi
echo ""

echo "▶ Levantando backend (Go) en ${BIND_HOST}:8080"
(cd "$ROOT_DIR/backend" && MC_HOST="$BIND_HOST" MC_ALLOWED_ORIGINS="$ORIGINS" go run .) &
BACKEND_PID=$!

echo "▶ Levantando frontend (SvelteKit) en ${BIND_HOST}:${UI_PORT}"
(cd "$ROOT_DIR/frontend" && pnpm dev --host "$BIND_HOST" --port "$UI_PORT" --strictPort) &
FRONTEND_PID=$!

wait "$BACKEND_PID" "$FRONTEND_PID"
