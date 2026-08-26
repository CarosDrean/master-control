import { getApiBase } from './api';
import { executeAction, type Action } from './actions';
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

const BASE_URL = getApiBase();

export async function getSystemActions(): Promise<SystemAction[]> {
    const res = await fetch(`${BASE_URL}/system/actions`);
    return res.json();
}

export async function getSystemVolume(): Promise<SystemVolume> {
    const res = await fetch(`${BASE_URL}/system/volume`);
    return res.json();
}

export async function getSystemBrightness(): Promise<SystemBrightness> {
    const res = await fetch(`${BASE_URL}/system/brightness`);
    return res.json();
}

export async function setSystemVolume(volume: number): Promise<void> {
    await fetch(`${BASE_URL}/system/volume`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ volume })
    });
}

export async function setSystemBrightness(brightness: number): Promise<void> {
    await fetch(`${BASE_URL}/system/brightness`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ brightness })
    });
}

export async function executeSystemActionWithFeedback(action: Action): Promise<void> {
    await executeAction(action.id);

    if (action.type !== "system") return;

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
}