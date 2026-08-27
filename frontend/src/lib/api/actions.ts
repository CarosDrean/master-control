import { getApiBase } from './api';

export interface Action {
    id: string;
    label: string;
    icon?: string;
    type: string;
    command: string;
    page?: string;
    position?: number;
}

const BASE_URL = getApiBase();

export async function getActions(): Promise<Action[]> {
    const res = await fetch(`${BASE_URL}/actions`);
    return res.json();
}

export async function getAction(id: string): Promise<Action> {
    const res = await fetch(`${BASE_URL}/actions/${id}`);
    return res.json();
}

export async function createAction(a: Action): Promise<{ success: boolean; message: string }> {
    const res = await fetch(`${BASE_URL}/actions`, {
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
    const res = await fetch(`${BASE_URL}/actions/${a.id}`, {
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
    await fetch(`${BASE_URL}/actions/${id}`, {
        method: "DELETE"
    });
}

export async function executeAction(id: string): Promise<void> {
    await fetch(`${BASE_URL}/actions/${id}/execute`, {
        method: "POST",
    });
}
