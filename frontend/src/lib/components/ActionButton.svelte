<script lang="ts">
    import { executeAction } from '$lib/api/actions';
    import { executeSystemActionWithFeedback } from '$lib/api/system';

    export let id: string;
    export let label: string;
    export let icon: string | undefined;
    export let type: string | undefined;
    export let command: string | undefined;
    export let onMenu: () => void = () => {};

    const isImage = !!icon && (icon.startsWith("http") || icon.startsWith("/") || /\.(png|jpe?g|svg|webp|gif|ico)$/i.test(icon));

    let pressTimer: ReturnType<typeof setTimeout> | null = null;
    let longPressed = false;

    function startPress(e: PointerEvent) {
        if (e.button !== 0) return;
        longPressed = false;
        pressTimer = setTimeout(() => {
            longPressed = true;
            onMenu();
        }, 600);
    }

    function endPress() {
        if (pressTimer) {
            clearTimeout(pressTimer);
            pressTimer = null;
        }
    }

    const handleClick = () => {
        if (longPressed) {
            longPressed = false;
            return;
        }
        if (type === "system" && command) {
            executeSystemActionWithFeedback({ id, label, icon: icon, type, command });
        } else {
            executeAction(id);
        }
    };

    const handleContext = (e: MouseEvent) => {
        e.preventDefault();
        onMenu();
    };
</script>

<style>
    .wrap {
        position: relative;
        width: 140px;
        height: 140px;
    }

    button.action {
        width: 140px;
        height: 140px;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        border-radius: 22px;
        border: 2px solid rgba(255,255,255,0.15);

        background: linear-gradient(180deg, #1f1f1f, #0e0e0e);
        box-shadow:
                0 6px 16px rgba(0,0,0,0.55),
                inset 0 0 0 3px rgba(255,255,255,0.07),
                inset 0 5px 14px rgba(255,255,255,0.03),
                inset 0 -5px 14px rgba(0,0,0,0.45);

        cursor: pointer;
        transition: 0.16s ease;
        touch-action: manipulation;
        user-select: none;
        -webkit-tap-highlight-color: transparent;
        -webkit-user-select: none;
    }

    button.action:active {
        transform: scale(0.95);
        box-shadow:
                0 2px 6px rgba(0,0,0,0.4),
                inset 0 0 0 3px rgba(255,255,255,0.05),
                inset 0 3px 6px rgba(255,255,255,0.01),
                inset 0 -3px 6px rgba(0,0,0,0.4);
    }

    .more {
        position: absolute;
        top: 8px;
        right: 8px;
        width: 28px;
        height: 28px;
        border-radius: 9px;
        border: 1px solid rgba(255,255,255,0.25);
        background: rgba(20,20,20,0.85);
        color: rgba(255,255,255,0.85);
        font-size: 1rem;
        line-height: 1;
        display: none;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        padding: 0;
        transition: 0.15s ease;
    }

    .more:hover {
        background: rgba(60,60,60,0.9);
    }

    .wrap:hover .more,
    .wrap:focus-within .more {
        display: flex;
    }

    img {
        width: 48px;
        height: 48px;
        margin-bottom: 6px;
        pointer-events: none;
    }

    .label {
        font-size: 0.85rem;
        opacity: 0.9;
        margin-top: 4px;
        color: rgba(255, 255, 255, 0.65);
        pointer-events: none;
    }

    .emoji {
        font-size: 34px;
        color: white;
        text-align: center;
        margin-bottom: 6px;
        pointer-events: none;
    }
</style>

<div class="wrap">
    <button
        class="action"
        type="button"
        on:click={handleClick}
        on:pointerdown={startPress}
        on:pointerup={endPress}
        on:pointerleave={endPress}
        on:pointercancel={endPress}
        on:contextmenu={handleContext}
    >
        {#if isImage}
            <img src={icon} alt={label} />
        {:else}
            <div class="emoji">{icon}</div>
        {/if}
        <div class="label">{label}</div>
    </button>

    <button class="more" type="button" aria-label="Opciones" on:click={onMenu}>⋯</button>
</div>

