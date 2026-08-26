export type HudData = {
    icon: string;
    label: string;
    category: "volume" | "brightness" | "media";
    value?: number;
    muted?: boolean;
};

type Listener = (data: HudData) => void;

let listener: Listener | null = null;

export function registerHud(cb: Listener): () => void {
    listener = cb;
    return () => {
        if (listener === cb) {
            listener = null;
        }
    };
}

export function showSystemHud(data: HudData): void {
    listener?.(data);
}