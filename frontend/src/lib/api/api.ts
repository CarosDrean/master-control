export const API_PORT = 8080;

const TOKEN_KEY = "master-control:token";

export function getApiBase(): string {
    if (typeof window === "undefined") {
        return `http://localhost:${API_PORT}/master-control`;
    }
    return `http://${window.location.hostname}:${API_PORT}/master-control`;
}

// ── Token ─────────────────────────────────────────────────────────────────
// El backend exige un token (MC_TOKEN). Llega una vez en la URL (`#token=…`,
// la que imprime dev.sh) o se pega a mano; se recuerda en localStorage.

let memoryToken = "";

type AuthListener = () => void;
const authListeners = new Set<AuthListener>();

/** Se dispara cuando el backend rechaza el token (401) o cuando falta. */
export function onAuthRequired(cb: AuthListener): () => void {
    authListeners.add(cb);
    return () => authListeners.delete(cb);
}

function emitAuthRequired() {
    for (const cb of authListeners) cb();
}

export function getToken(): string {
    if (memoryToken) return memoryToken;
    try {
        memoryToken = localStorage.getItem(TOKEN_KEY) ?? "";
    } catch {
        // almacenamiento bloqueado: el token vive solo en memoria
    }
    return memoryToken;
}

export function setToken(token: string): void {
    memoryToken = token.trim();
    try {
        if (memoryToken) localStorage.setItem(TOKEN_KEY, memoryToken);
        else localStorage.removeItem(TOKEN_KEY);
    } catch {
        // almacenamiento bloqueado: el token vive solo en memoria
    }
}

/** Lee `#token=…` de la URL, lo guarda y lo quita de la barra de direcciones. */
export function adoptTokenFromUrl(): void {
    if (typeof window === "undefined") return;
    const hash = new URLSearchParams(window.location.hash.slice(1));
    const token = hash.get("token");
    if (!token) return;
    setToken(token);
    hash.delete("token");
    const rest = hash.toString();
    history.replaceState(null, "", window.location.pathname + window.location.search + (rest ? `#${rest}` : ""));
}

// ── Cliente ───────────────────────────────────────────────────────────────

/** fetch contra la API con el token en `Authorization`. Un 401 pide el token de nuevo. */
export async function apiFetch(path: string, init: RequestInit = {}): Promise<Response> {
    const headers = new Headers(init.headers);
    const token = getToken();
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const res = await fetch(`${getApiBase()}${path}`, { ...init, headers });
    if (res.status === 401) emitAuthRequired();
    return res;
}

/** Comprueba el token contra un endpoint protegido. */
export async function verifyToken(): Promise<boolean> {
    try {
        const res = await apiFetch("/action-types");
        return res.ok;
    } catch {
        return false;
    }
}

/**
 * URL para recursos que el navegador pide sin cabeceras (<img src>).
 * Solo añade el token a URLs de esta API; nunca se persiste en la acción.
 */
export function authUrl(url: string | undefined): string | undefined {
    if (!url) return url;
    const token = getToken();
    if (!token) return url;
    const base = getApiBase();
    // Los iconos guardados pueden tener otro host (p. ej. la IP de la LAN): se reescriben al actual.
    const marker = "/master-control/";
    const at = url.indexOf(marker);
    if (!url.startsWith(base) && !(url.startsWith("http") && at >= 0 && url.includes(`:${API_PORT}`))) {
        return url;
    }
    const local = url.startsWith(base) ? url : base + "/" + url.slice(at + marker.length);
    const sep = local.includes("?") ? "&" : "?";
    return `${local}${sep}token=${encodeURIComponent(token)}`;
}
