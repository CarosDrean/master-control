/** "" o inválido → undefined (timeout por defecto del backend). */
export function parseTimeout(raw: string): number | undefined {
    const n = parseInt(String(raw ?? "").trim(), 10);
    return Number.isFinite(n) && n > 0 ? n : undefined;
}
