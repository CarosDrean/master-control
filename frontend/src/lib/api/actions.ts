import { apiFetch } from './api';

export interface Action {
    id: string;
    label: string;
    icon?: string;
    type: string;
    command: string;
    page?: string;
    position?: number;
    /** Pide confirmación antes de ejecutar. */
    confirm?: boolean;
    /** Muestra la salida del comando al terminar. */
    show_output?: boolean;
    /** Timeout en segundos (0/vacío = por defecto; el backend lo recorta a MC_MAX_TIMEOUT). */
    timeout_seconds?: number;
}

export interface ExecutionResult {
    ok: boolean;
    exit_code: number;
    stdout: string;
    stderr: string;
    duration_ms: number;
    timed_out: boolean;
    truncated: boolean;
}

export async function getActions(): Promise<Action[]> {
    const res = await apiFetch(`/actions`);
    return res.json();
}

export async function getAction(id: string): Promise<Action> {
    const res = await apiFetch(`/actions/${id}`);
    return res.json();
}

export async function createAction(a: Action): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/actions`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(a)
    });

    if (!res.ok) {
        return {
            success: false,
            message: `Error ${res.status}: ${res.statusText}`
        };
    }

    try {
        const data = await res.json();
        return {
            success: data.success ?? true,
            message: data.message ?? "Guardado correctamente"
        };
    } catch {
        return {
            success: true,
            message: "Guardado correctamente"
        };
    }
}

export async function updateAction(a: Action): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/actions/${a.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(a)
    });

    if (!res.ok) {
        return {
            success: false,
            message: `Error ${res.status}: ${res.statusText}`
        };
    }

    try {
        const data = await res.json();
        return {
            success: data.success ?? true,
            message: data.message ?? "Guardado correctamente"
        };
    } catch {
        return {
            success: true,
            message: "Guardado correctamente"
        };
    }
}


export async function deleteAction(id: string): Promise<void> {
    await apiFetch(`/actions/${id}`, {
        method: "DELETE"
    });
}

export async function executeAction(id: string): Promise<ExecutionResult> {
    let res: Response;
    try {
        res = await apiFetch(`/actions/${id}/execute`, { method: "POST" });
    } catch {
        return failure("No se pudo conectar con el servidor");
    }

    let data: unknown = null;
    try {
        data = await res.json();
    } catch { /* cuerpo vacío o no JSON */ }

    if (!res.ok) {
        const message = (data as { error?: string } | null)?.error ?? `Error ${res.status}: ${res.statusText}`;
        return failure(message);
    }

    return data as ExecutionResult;
}

function failure(message: string): ExecutionResult {
    return { ok: false, exit_code: -1, stdout: "", stderr: message, duration_ms: 0, timed_out: false, truncated: false };
}
