import { apiFetch } from './api';

export interface Page {
    id: string;
    name: string;
    icon?: string;
    order: number;
    columns?: number;
    rows?: number;
    /** App de la Mac vinculada (nombre del bundle sin ".app"): la UI salta aquí cuando pasa al frente. */
    app?: string;
}

export async function getPages(): Promise<Page[]> {
    const res = await apiFetch(`/pages`);
    return res.json();
}

export async function getPage(id: string): Promise<Page> {
    const res = await apiFetch(`/pages/${id}`);
    return res.json();
}

export async function createPage(p: Page): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/pages`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(p)
    });

    if (!res.ok) {
        return { success: false, message: `Error ${res.status}: ${res.statusText}` };
    }

    return { success: true, message: "Página guardada" };
}

export async function updatePage(p: Page): Promise<{ success: boolean; message: string }> {
    const res = await apiFetch(`/pages/${p.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(p)
    });

    if (!res.ok) {
        return { success: false, message: `Error ${res.status}: ${res.statusText}` };
    }

    return { success: true, message: "Página actualizada" };
}

export async function deletePage(id: string): Promise<void> {
    await apiFetch(`/pages/${id}`, {
        method: "DELETE"
    });
}