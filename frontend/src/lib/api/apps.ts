import { apiFetch, getApiBase } from './api';

export interface App {
    name: string;
    path: string;
}

const BASE_URL = getApiBase();

export function appIconUrl(path: string): string {
    return `${BASE_URL}/apps/icon?path=${encodeURIComponent(path)}`;
}

export function appActionId(name: string): string {
    const slug = name.toLowerCase().replace(/[^a-z0-9]+/g, "_").replace(/^_+|_+$/g, "");
    return `open_${slug}_app`;
}

export async function getApps(): Promise<App[]> {
    const res = await apiFetch(`/apps`);
    return res.json();
}