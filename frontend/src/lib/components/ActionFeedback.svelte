<script lang="ts">
    import { onMount } from 'svelte';
    import { registerFeedback, type RunFeedback } from '$lib/run';
    import type { Action } from '$lib/api/actions';

    let confirming = $state<Action | null>(null);
    let resolveConfirm: ((ok: boolean) => void) | null = null;

    let running = $state<string[]>([]);
    let panel = $state<RunFeedback | null>(null);
    let toast = $state<RunFeedback | null>(null);
    let toastTimer: ReturnType<typeof setTimeout> | null = null;

    function answer(ok: boolean) {
        resolveConfirm?.(ok);
        resolveConfirm = null;
        confirming = null;
    }

    onMount(() => registerFeedback({
        confirm(action) {
            // Una confirmación pendiente se cancela si llega otra.
            resolveConfirm?.(false);
            confirming = action;
            return new Promise<boolean>((resolve) => (resolveConfirm = resolve));
        },
        running(action, on) {
            running = on ? [...running, action.id] : running.filter((id) => id !== action.id);
        },
        show(feedback) {
            if (feedback.expanded) {
                panel = feedback;
                return;
            }
            toast = feedback;
            if (toastTimer) clearTimeout(toastTimer);
            toastTimer = setTimeout(() => (toast = null), 1600);
        },
    }));

    function onKey(e: KeyboardEvent) {
        if (e.key !== 'Escape') return;
        if (confirming) answer(false);
        else panel = null;
    }

    const status = (f: RunFeedback) =>
        f.result.timed_out ? 'Tiempo agotado' : f.result.ok ? 'Correcto' : `Falló (código ${f.result.exit_code})`;

    const seconds = (ms: number) => (ms / 1000).toFixed(ms < 10000 ? 1 : 0) + ' s';

    async function copy(text: string) {
        try {
            await navigator.clipboard.writeText(text);
        } catch { /* sin permiso de portapapeles */ }
    }
</script>

<svelte:window onkeydown={onKey} />

{#if running.length > 0}
    <div class="running" aria-live="polite">
        <span class="spinner"></span> Ejecutando…
    </div>
{/if}

{#if toast}
    <div class="toast" class:fail={!toast.result.ok} aria-live="polite">
        <span>{toast.result.ok ? '✓' : '✕'}</span> {toast.action.label}
    </div>
{/if}

{#if confirming}
    <div class="overlay" role="presentation" onclick={() => answer(false)}>
        <div class="dialog" role="alertdialog" aria-modal="true" aria-labelledby="confirm-title" tabindex="-1" onclick={(e) => e.stopPropagation()} onkeydown={() => {}}>
            <h3 id="confirm-title">¿Ejecutar «{confirming.label}»?</h3>
            {#if confirming.command}
                <pre class="cmd">{confirming.command}</pre>
            {/if}
            <div class="buttons">
                <button class="btn" onclick={() => answer(false)}>Cancelar</button>
                <button class="btn primary" onclick={() => answer(true)}>Ejecutar</button>
            </div>
        </div>
    </div>
{/if}

{#if panel}
    {@const r = panel.result}
    <div class="overlay" role="presentation" onclick={() => (panel = null)}>
        <div class="dialog wide" role="dialog" aria-modal="true" aria-labelledby="out-title" tabindex="-1" onclick={(e) => e.stopPropagation()} onkeydown={() => {}}>
            <header>
                <h3 id="out-title">{panel.action.label}</h3>
                <span class="badge" class:fail={!r.ok}>{status(panel)}</span>
                <span class="muted">{seconds(r.duration_ms)}</span>
            </header>

            {#if r.stdout}
                <pre class="out">{r.stdout}</pre>
            {/if}
            {#if r.stderr}
                <pre class="out err">{r.stderr}</pre>
            {/if}
            {#if !r.stdout && !r.stderr}
                <p class="muted">Sin salida.</p>
            {/if}
            {#if r.truncated}
                <p class="muted">Salida recortada a 16 KB por stream.</p>
            {/if}

            <div class="buttons">
                {#if r.stdout || r.stderr}
                    <button class="btn" onclick={() => copy([r.stdout, r.stderr].filter(Boolean).join('\n'))}>Copiar</button>
                {/if}
                <button class="btn primary" onclick={() => (panel = null)}>Cerrar</button>
            </div>
        </div>
    </div>
{/if}

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.6);
        z-index: 80;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 24px;
    }

    .dialog {
        width: 100%;
        max-width: 420px;
        max-height: 85vh;
        overflow: auto;
        background: #161616;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 18px;
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
        padding: 20px;
        color: white;
    }

    .dialog.wide {
        max-width: 760px;
    }

    header {
        display: flex;
        align-items: center;
        gap: 10px;
        flex-wrap: wrap;
        margin-bottom: 14px;
    }

    h3 {
        margin: 0;
        font-size: 1.1rem;
    }

    .dialog > h3 {
        margin-bottom: 14px;
    }

    pre {
        margin: 0 0 12px;
        padding: 12px;
        border-radius: 10px;
        background: #0c0c0c;
        border: 1px solid rgba(255, 255, 255, 0.08);
        font-size: 0.85rem;
        line-height: 1.45;
        white-space: pre-wrap;
        word-break: break-word;
        max-height: 50vh;
        overflow: auto;
    }

    .cmd {
        max-height: 140px;
        color: rgba(255, 255, 255, 0.7);
    }

    .err {
        color: #ffb3b3;
    }

    .badge {
        font-size: 0.78rem;
        padding: 3px 10px;
        border-radius: 999px;
        background: rgba(90, 200, 120, 0.18);
        color: #9be8b4;
    }

    .badge.fail {
        background: rgba(255, 100, 100, 0.18);
        color: #ffb3b3;
    }

    .muted {
        color: rgba(255, 255, 255, 0.5);
        font-size: 0.85rem;
    }

    .buttons {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        margin-top: 8px;
    }

    .btn {
        min-height: 48px;
        min-width: 110px;
        padding: 0 18px;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.15);
        background: #1e1e1e;
        color: white;
        font-size: 1rem;
        cursor: pointer;
        touch-action: manipulation;
    }

    .btn.primary {
        background: #2a5bd7;
        border-color: #2a5bd7;
    }

    .toast,
    .running {
        position: fixed;
        left: 50%;
        transform: translateX(-50%);
        z-index: 90;
        padding: 10px 18px;
        border-radius: 999px;
        background: rgba(30, 30, 30, 0.95);
        border: 1px solid rgba(255, 255, 255, 0.12);
        color: white;
        font-size: 0.95rem;
        display: flex;
        align-items: center;
        gap: 8px;
        pointer-events: none;
    }

    .toast {
        bottom: 28px;
    }

    .toast.fail {
        border-color: rgba(255, 100, 100, 0.5);
    }

    .running {
        top: 72px;
    }

    .spinner {
        width: 14px;
        height: 14px;
        border-radius: 50%;
        border: 2px solid rgba(255, 255, 255, 0.25);
        border-top-color: white;
        animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }
</style>
