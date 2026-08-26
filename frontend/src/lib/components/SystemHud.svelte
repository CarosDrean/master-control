<script lang="ts">
    import { onMount } from 'svelte';
    import { registerHud, type HudData } from '$lib/hud';
    import { setSystemVolume, setSystemBrightness } from '$lib/api/system';

    let hud = $state<HudData | null>(null);
    let visible = $state(false);
    let timer: ReturnType<typeof setTimeout> | null = null;

    let applying = false;
    let pendingValue: number | null = null;

    function onHud(data: HudData) {
        hud = data;
        visible = true;
        if (timer) clearTimeout(timer);
        timer = setTimeout(() => (visible = false), 1300);
    }

    onMount(() => registerHud(onHud));

    function resetHideTimer() {
        if (timer) clearTimeout(timer);
        timer = setTimeout(() => (visible = false), 1300);
    }

    function barPointerDown(e: PointerEvent) {
        if (timer) clearTimeout(timer);
        const el = e.currentTarget as HTMLElement;
        el.setPointerCapture(e.pointerId);
        applyFromPointer(e);
    }

    function barPointerMove(e: PointerEvent) {
        const el = e.currentTarget as HTMLElement;
        if (el.hasPointerCapture(e.pointerId)) {
            applyFromPointer(e);
        }
    }

    function barPointerUp(e: PointerEvent) {
        const el = e.currentTarget as HTMLElement;
        if (el.hasPointerCapture(e.pointerId)) {
            el.releasePointerCapture(e.pointerId);
        }
        resetHideTimer();
    }

    function applyFromPointer(e: PointerEvent) {
        if (!hud) return;

        const el = e.currentTarget as HTMLElement;
        const rect = el.getBoundingClientRect();
        let pct = ((e.clientX - rect.left) / rect.width) * 100;
        pct = Math.max(0, Math.min(100, Math.round(pct)));

        hud.value = pct;
        if (hud.category === "volume") {
            hud.muted = false;
        }

        pendingValue = pct;
        if (applying) return;
        applying = true;

        void (async () => {
            while (pendingValue !== null) {
                const value = pendingValue;
                pendingValue = null;
                try {
                    if (hud?.category === "volume") {
                        await setSystemVolume(value);
                    } else {
                        await setSystemBrightness(value / 100);
                    }
                } catch {
                    // ignora errores transitorios de red
                }
            }
            applying = false;
        })();
    }
</script>

<style>
    .hud-overlay {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background: rgba(18, 18, 20, 0.94);
        border: 1px solid rgba(255, 255, 255, 0.12);
        border-radius: 22px;
        padding: 26px 44px;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
        box-shadow: 0 18px 60px rgba(0, 0, 0, 0.6);
        opacity: 0;
        transition: opacity 0.25s ease;
        pointer-events: none;
        z-index: 9999;
    }

    .hud-overlay.visible {
        opacity: 1;
        pointer-events: auto;
    }

    .hud-icon {
        font-size: 42px;
        line-height: 1;
    }

    .hud-label {
        font-size: 0.95rem;
        opacity: 0.9;
    }

    .hud-bar {
        position: relative;
        width: 220px;
        height: 16px;
        border-radius: 8px;
        background: rgba(255, 255, 255, 0.12);
        overflow: hidden;
        margin-top: 6px;
        cursor: pointer;
        touch-action: none;
        user-select: none;
        -webkit-user-select: none;
    }

    .hud-fill {
        height: 100%;
        background: linear-gradient(90deg, #5ab3ff, #8a8fff);
        border-radius: 8px;
        transition: width 0.08s linear;
    }

    .hud-thumb {
        position: absolute;
        top: 50%;
        transform: translate(-50%, -50%);
        width: 18px;
        height: 18px;
        border-radius: 50%;
        background: #fff;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.5);
    }

    .hud-value {
        font-size: 0.85rem;
        opacity: 0.75;
    }
</style>

{#if hud}
    <div class="hud-overlay" class:visible>
        <div class="hud-icon">{hud.icon}</div>
        <div class="hud-label">{hud.label}</div>
        {#if hud.category === "volume" || hud.category === "brightness"}
            <div
                class="hud-bar"
                role="slider"
                tabindex="0"
                aria-valuemin="0"
                aria-valuemax="100"
                aria-valuenow={hud.value ?? 0}
                onpointerdown={barPointerDown}
                onpointermove={barPointerMove}
                onpointerup={barPointerUp}
                onpointercancel={barPointerUp}
            >
                <div class="hud-fill" style:width="{hud.value ?? 0}%"></div>
                <div class="hud-thumb" style:left="{hud.value ?? 0}%"></div>
            </div>
            <div class="hud-value">
                {Math.round(hud.value ?? 0)}%{hud.category === "volume" && hud.muted ? " · silenciado" : ""}
            </div>
        {/if}
    </div>
{/if}