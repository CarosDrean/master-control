# Master Control

Monorepo de Master Control: panel para controlar tu Mac (acciones, hotkeys y páginas) desde el navegador.

## Estructura

```
.
├── backend/    # API en Go (Gin + bbolt), expuesta en :8080/master-control
├── frontend/   # UI en SvelteKit + Vite (pnpm)
└── bbolt.db    # Base de datos (creada en la raíz al ejecutar el backend)
```

## Requisitos

- Go 1.25+
- Node.js + pnpm

## Desarrollo

Levanta ambos servicios desde la raíz:

```sh
./.dev.sh
```

(equivalente a correr `go run .` en `backend/` y `pnpm dev` en `frontend/`)

## Backend

```sh
cd backend
go run .
```

Servidor en `http://localhost:8080/master-control`. Configuración vía `backend/.env` (`GIN_MODE`, `APP_NAME`, `BUCKET_NAME`). Ejemplos de requests en `backend/api.http`.

## Frontend

```sh
cd frontend
pnpm install
pnpm dev
```

La UI apunta a la API del backend en `http://<host>:8080/master-control` (ver `frontend/src/lib/api/api.ts`).

## Notas

- `bbolt.db` se crea en la raíz del repo (relativo al cwd del backend).
- La base de datos y `.env` están en `.gitignore`.