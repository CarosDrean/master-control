# AGENTS.md

Master Control: macOS control panel (trigger actions/hotkeys/apps on the Mac from the browser). Monorepo — Go backend + SvelteKit frontend. All UI strings and code comments are in Spanish.

## Layout

- `backend/` — Go API (Gin + bbolt). Route prefix `/master-control`, `127.0.0.1:8080` by default (`MC_HOST`/`MC_PORT`). Module `master-control`, Go 1.25.
- `frontend/` — SvelteKit 5 + Vite + TypeScript, managed with **pnpm** (has `pnpm-workspace.yaml`). No root JS workspace.

## Dev commands (run from each dir; there is no root package.json)

```sh
./dev.sh [--tailscale|--lan]    # both services together (go run . + pnpm dev); default = loopback only
cd backend && go run .          # backend on http://localhost:8080/master-control
cd frontend && pnpm install && pnpm dev   # SvelteKit
```

Verification:
- Backend: `go build ./...` / `go vet ./...` / `go test ./...` (config, auth middleware, action execution)
- Frontend: `pnpm check` (svelte-check typecheck). `pnpm build` = `vite build`.

## Backend gotchas

- **Run it from `backend/`.** `godotenv.Load()` reads `.env` and bbolt opens `bbolt.db` relative to the process cwd. The repo-root `bbolt.db` is stale state; don't touch it.
- `backend/.env` is gitignored (no committed `.env.example`). Missing `.env` is not fatal: `GIN_MODE` falls back to gin debug mode and `APP_NAME` falls back to empty → routes are served at `/` instead of `/master-control`. Keep `APP_NAME=master-control`.
- macOS-only features: `infrastructure/hotkey` uses **cgo** (Carbon/ApplicationServices, `//go:build darwin`); non-darwin builds return `ErrNotSupported`. `POST /hotkey/capture` blocks up to 15s waiting for a keypress and requires macOS Accessibility permission.
- Actions are executed as shell commands / AppleScript (`osascript` → System Events). App listing/icon endpoints scan `/Applications` and shell out to `sips`.
- Layered architecture: `bootstrap/` wires use-cases (`domain/`) to repos & handlers (`infrastructure/`). Repos persist JSON into bbolt buckets directly; `action-types` and `pages` buckets are seeded on startup.
- **Security model (keep it):** every route except `/health` requires `MC_TOKEN` (`Authorization: Bearer`, or `?token=` for `<img>`), compared in constant time (`infrastructure/middleware/auth.go`). No `MC_TOKEN` → an ephemeral one is generated and logged; the server never runs open. The request logger strips the query string so the token never lands in logs. CORS is limited to `MC_ALLOWED_ORIGINS`. Never bind `0.0.0.0` by default; `dev.sh --lan` is the explicit opt-in.
- `POST /actions/:id/execute` returns a `model.ExecutionResult` (stdout/stderr capped at 16 KB each). A non-zero exit or timeout is `200` with `ok:false`; `5xx` only when the action cannot start. Per-action `timeout_seconds` is capped by `MC_MAX_TIMEOUT`; the server `WriteTimeout` derives from it.
- `api.http` has sample requests (VS Code REST Client) using `@base`/`@token` variables.

## Frontend gotchas

- API base is hardcoded in `src/lib/api/api.ts`: `http://<host>:8080/master-control`. No Vite proxy; relies on backend CORS. Change the port in one place only.
- **Always call the API through `apiFetch`** (adds the token; a 401 re-opens the token screen, `TokenGate.svelte`). Wrap API image URLs in `authUrl()` at render time — never persist a token into an action's `icon`.
- Run actions through `runAction()` (`src/lib/run.ts`): it handles `confirm`, the running indicator and the result (toast or output panel, `ActionFeedback.svelte`, mounted in the layout).
- Per-resource API clients live in `src/lib/api/` (one file per resource); components in `src/lib/components/`; routes under `src/routes/{actions,pages}`.
- Svelte 5 runes-style; Svelte 5 syntax applies (no legacy stores unless already used).