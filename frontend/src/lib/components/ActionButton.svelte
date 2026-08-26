<script lang="ts">
    import { executeAction } from '$lib/api/actions';
    import { executeSystemActionWithFeedback } from '$lib/api/system';

    export let id: string;
    export let label: string;
    export let icon: string | undefined;
    export let type: string | undefined;
    export let command: string | undefined;

    const isImage = !!icon && (icon.startsWith("http") || icon.startsWith("/") || /\.(png|jpe?g|svg|webp|gif|ico)$/i.test(icon));

    const handleClick = () => {
        if (type === "system" && command) {
            executeSystemActionWithFeedback({ id, label, icon: icon, type, command });
        } else {
            executeAction(id);
        }
    };
</script>

<style>
    button {
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

    button:active {
        transform: scale(0.95);
        box-shadow:
                0 2px 6px rgba(0,0,0,0.4),
                inset 0 0 0 3px rgba(255,255,255,0.05),
                inset 0 3px 6px rgba(255,255,255,0.01),
                inset 0 -3px 6px rgba(0,0,0,0.4);
    }

    img {
        width: 48px;
        height: 48px;
        margin-bottom: 6px;
    }

    .label {
        font-size: 0.85rem;
        opacity: 0.9;
        margin-top: 4px;
        color: rgba(255, 255, 255, 0.65);
    }

    .emoji {
        font-size: 34px;
        color: white;
        text-align: center;
        margin-bottom: 6px;
    }
</style>


<button type="button" on:click={handleClick}>
    {#if isImage}
        <img src={icon} alt={label} />
    {:else}
        <div class="emoji">{icon}</div>
    {/if}
    <div class="label">{label}</div>
</button>

