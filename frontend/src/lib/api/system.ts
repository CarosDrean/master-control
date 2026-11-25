import { apiFetch } from './api';
import { executeAction, type Action, type ExecutionResult } from './actions';
import { showSystemHud } from '$lib/hud';

export interface SystemAction {
    id: string;
    label: string;
    icon: string;
    command: string;
    category: string;
}

export interface SystemVolume {
    volume: number;
    muted: boolean;
}

export interface SystemBrightness {
    brightness: number;
}

export async function getSystemActions(): Promise<SystemAction[]> {
    const res = await apiFetch(`/system/actions`);
    return res.json();
}

export async function getSystemVolume(): Promise<SystemVolume> {
    const res = await apiFetch(`/system/volume`);
    return res.json();
}

export async function getSystemBrightness(): Promise<SystemBrightness> {
    const res = await apiFetch(`/system/brightness`);
    return res.json();
}

export async function setSystemVolume(volume: number): Promise<void> {
    await apiFetch(`/system/volume`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ volume })
    });
}

export async function setSystemBrightness(brightness: number): Promise<void> {
    await apiFetch(`/system/brightness`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ brightness })
    });
}

export async function executeSystemActionWithFeedback(action: Action): Promise<ExecutionResult> {
    const result = await executeAction(action.id);

    if (action.type !== "system" || !result.ok) return result;

    const category: "volume" | "brightness" | "media" = action.command.startsWith("volume")
        ? "volume"
        : action.command.startsWith("brightness")
            ? "brightness"
            : "media";

    let value: number | undefined;
    let muted: boolean | undefined;

    try {
        if (category === "volume") {
            const v = await getSystemVolume();
            value = v.volume;
            muted = v.muted;
        } else if (category === "brightness") {
            const b = await getSystemBrightness();
            value = Math.round(b.brightness * 100);
        }
    } catch {
        // sin valor disponible, se muestra el HUD igual
    }

    showSystemHud({ icon: action.icon ?? "⚙️", label: action.label, category, value, muted });

    return result;
}

export type ScrollResult = { ok: true } | { ok: false; status: number; error: string };

/**
 * Desplaza la ventana de la Mac que está bajo el puntero, en píxeles.
 * Positivo = hacia arriba / izquierda (convención de la rueda de CoreGraphics).
 */
export async function scrollBy(dx: number, dy: number): Promise<ScrollResult> {
    const res = await apiFetch(`/system/scroll`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ dx: Math.round(dx), dy: Math.round(dy) })
    });
    if (res.ok) return { ok: true };
    let error = `HTTP ${res.status}`;
    try {
        error = (await res.json()).error ?? error;
    } catch {
        // cuerpo vacío o no JSON
    }
    return { ok: false, status: res.status, error };
}
