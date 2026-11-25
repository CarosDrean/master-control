import { apiFetch } from './api';

export interface ActionType {
    id: string;
    key: string;
    name: string;
    description: string;
    command_template: string;
}

export async function getActionTypes(): Promise<ActionType[]> {
    const res = await apiFetch(`/action-types`);
    return res.json();
}

export async function getActionType(id: string): Promise<ActionType> {
    const res = await apiFetch(`/action-types/${id}`);
    return res.json();
}

export async function createActionType(a: ActionType): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/action-types`, {
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

export async function updateActionType(a: ActionType): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/action-types/${a.id}`, {
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


export async function deleteActionType(id: string): Promise<void> {
    await apiFetch(`/action-types/${id}`, {
        method: "DELETE"
    });
}

