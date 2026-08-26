# AGENTS.md

Master Control: macOS control panel (trigger actions/hotkeys/apps on the Mac from the browser). Monorepo — Go backend + SvelteKit frontend. All UI strings and code comments are in Spanish.

## Layout

- `backend/` — Go API (Gin + bbolt). Route prefix `/master-control`, port `:8080`. Module `master-control`, Go 1.25.
- `frontend/` — SvelteKit 5 + Vite + TypeScript, managed with **pnpm** (has `pnpm-workspace.yaml`). No root JS workspace.

## Dev commands (run from each dir; there is no root package.json)

```sh
./.dev.sh                       # both services together (go run . + pnpm dev)
cd backend && go run .          # backend on http://localhost:8080/master-control
cd frontend && pnpm install && pnpm dev   # SvelteKit
```

Verification: no tests or linters exist anywhere.
- Backend: `go build ./...` / `go vet ./...`
- Frontend: `pnpm check` (svelte-check typecheck). `pnpm build` = `vite build`.

## Backend gotchas

- **Run it from `backend/`.** `godotenv.Load()` reads `.env` and bbolt opens `bbolt.db` relative to the process cwd. The repo-root `bbolt.db` is stale state; don't touch it.
- `backend/.env` is gitignored (no committed `.env.example`). Missing `.env` is not fatal: `GIN_MODE` falls back to gin debug mode and `APP_NAME` falls back to empty → routes are served at `/` instead of `/master-control`. Keep `APP_NAME=master-control`.
- macOS-only features: `infrastructure/hotkey` uses **cgo** (Carbon/ApplicationServices, `//go:build darwin`); non-darwin builds return `ErrNotSupported`. `POST /hotkey/capture` blocks up to 15s waiting for a keypress and requires macOS Accessibility permission.
- Actions are executed as shell commands / AppleScript (`osascript` → System Events). App listing/icon endpoints scan `/Applications` and shell out to `sips`.
- Layered architecture: `bootstrap/` wires use-cases (`domain/`) to repos & handlers (`infrastructure/`). Repos persist JSON into bbolt buckets directly; `action-types` and `pages` buckets are seeded on startup.
- `api.http` has sample requests (VS Code REST Client) pointing at `0.0.0.0:8080`.

## Frontend gotchas

- API base is hardcoded in `src/lib/api/api.ts`: `http://<host>:8080/master-control`. No Vite proxy; relies on backend CORS (all origins). Change the port in one place only.
- Per-resource API clients live in `src/lib/api/` (one file per resource); components in `src/lib/components/`; routes under `src/routes/{actions,pages}`.
- Svelte 5 runes-style; Svelte 5 syntax applies (no legacy stores unless already used).