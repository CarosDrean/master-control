import { getApiBase } from './api';

export interface Hotkey {
    keycode: number;
    modifiers: string[];
    command: string;
}

export interface HotkeyCaptureResult {
    success: boolean;
    hotkey?: Hotkey;
    message?: string;
    status?: number;
}

const BASE_URL = getApiBase();

export async function captureHotkey(): Promise<HotkeyCaptureResult> {
    try {
        const res = await fetch(`${BASE_URL}/hotkey/capture`, {
            method: "POST"
        });

        if (!res.ok) {
            let message = `Error ${res.status}: ${res.statusText}`;
            try {
                const data = await res.json();
                if (data.error) message = data.error;
            } catch { /* ignore */ }

            return { success: false, message, status: res.status };
        }

        const hotkey: Hotkey = await res.json();
        return { success: true, hotkey };
    } catch (e) {
        return {
            success: false,
            message: "No se pudo conectar con el servidor"
        };
    }
}