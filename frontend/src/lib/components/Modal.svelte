<script lang="ts">
    export let title: string;
    export let onClose: () => void = () => {};
    export let width = 600;
    export let maxHeight = "88vh";
</script>

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.6);
        z-index: 60;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 24px;
    }

    .modal {
        width: 100%;
        max-width: 600px;
        max-height: 88vh;
        display: flex;
        flex-direction: column;
        background: #161616;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 18px;
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
        overflow: hidden;
    }

    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 18px 24px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    }

    .header h3 {
        margin: 0;
        font-size: 1.05rem;
    }

    .close {
        background: none;
        border: none;
        color: #aaa;
        font-size: 1.2rem;
        cursor: pointer;
        padding: 4px;
    }

    .body {
        overflow-y: auto;
        padding: 0 24px 24px;
    }
</style>

<div
    class="overlay"
    role="presentation"
    on:click={(e) => { if (e.target === e.currentTarget) onClose(); }}
    on:keydown={(e) => { if (e.key === 'Escape') onClose(); }}
>
    <div class="modal" style="max-width: {width}px; max-height: {maxHeight}" role="dialog" aria-label={title}>
        <div class="header">
            <h3>{title}</h3>
            <button class="close" on:click={() => onClose()} aria-label="Cerrar">✕</button>
        </div>
        <div class="body">
            <slot />
        </div>
    </div>
</div>