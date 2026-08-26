import { getApiBase } from './api';

export interface Page {
    id: string;
    name: string;
    icon?: string;
    order: number;
}

const BASE_URL = getApiBase();

export async function getPages(): Promise<Page[]> {
    const res = await fetch(`${BASE_URL}/pages`);
    return res.json();
}

export async function getPage(id: string): Promise<Page> {
    const res = await fetch(`${BASE_URL}/pages/${id}`);
    return res.json();
}

export async function createPage(p: Page): Promise<{ success: boolean; message: string }> {
    const res = await fetch(`${BASE_URL}/pages`, {
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
    const res = await fetch(`${BASE_URL}/pages/${p.id}`, {
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
    await fetch(`${BASE_URL}/pages/${id}`, {
        method: "DELETE"
    });
}