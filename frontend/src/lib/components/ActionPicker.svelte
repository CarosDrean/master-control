<script lang="ts">
    import type { Action } from '$lib/api/actions';

    export let actions: Action[] = [];
    export let pages: { id: string; name: string; icon?: string }[] = [];
    export let onSelect: (a: Action) => void = () => {};
    export let onNew: () => void = () => {};
    export let onClose: () => void = () => {};
    let search = "";

    $: filtered = search
        ? actions.filter(a => a.label.toLowerCase().includes(search.toLowerCase()))
        : actions;

    const pageName = (id?: string) => {
        const p = pages.find(x => x.id === id);
        return p ? `${p.icon ?? ""} ${p.name}`.trim() : id ?? "general";
    };

    const isImage = (icon?: string) =>
        !!icon && (icon.startsWith("http") || icon.startsWith("/") || /\.(png|jpe?g|svg|webp|gif|ico)$/i.test(icon));
</script>

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.6);
        z-index: 50;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 24px;
    }

    .modal {
        width: 100%;
        max-width: 460px;
        max-height: 70vh;
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
        padding: 18px 20px;
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

    .search {
        margin: 14px 20px;
        background: #1f1f1f;
        border: 1px solid rgba(255, 255, 255, 0.12);
        border-radius: 10px;
        padding: 10px 14px;
        color: white;
        font-size: 0.95rem;
        outline: none;
    }

    .list {
        flex: 1;
        overflow-y: auto;
        padding: 0 10px 10px;
    }

    .item {
        display: flex;
        align-items: center;
        gap: 12px;
        width: 100%;
        background: none;
        border: none;
        border-radius: 10px;
        padding: 10px 12px;
        color: white;
        font-size: 0.95rem;
        text-align: left;
        cursor: pointer;
    }

    .item:hover {
        background: rgba(255, 255, 255, 0.06);
    }

    .item img {
        width: 30px;
        height: 30px;
        border-radius: 6px;
        flex-shrink: 0;
    }

    .item .emoji {
        width: 30px;
        text-align: center;
        font-size: 1.2rem;
        flex-shrink: 0;
    }

    .item .label {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .item .page-badge {
        font-size: 0.75rem;
        opacity: 0.55;
        flex-shrink: 0;
    }

    .empty {
        text-align: center;
        opacity: 0.5;
        padding: 24px 0;
    }

    .new {
        display: block;
        width: calc(100% - 40px);
        margin: 0 20px 18px;
        background: #1e1e1e;
        border: 1px solid rgba(255, 255, 255, 0.12);
        border-radius: 12px;
        padding: 12px;
        font-size: 0.95rem;
        color: white;
        cursor: pointer;
    }

    .new:hover {
        background: #292929;
    }
</style>

<div class="overlay" role="presentation" on:click={(e) => { if (e.target === e.currentTarget) onClose(); }}>
    <div class="modal" role="dialog" aria-label="Seleccionar acción">
        <div class="header">
            <h3>Seleccionar acción</h3>
            <button class="close" on:click={() => onClose()} aria-label="Cerrar">✕</button>
        </div>

        <input class="search" bind:value={search} placeholder="Buscar acción…" />

        <div class="list">
            {#if filtered.length === 0}
                <div class="empty">No hay acciones sin asignar.</div>
            {:else}
                {#each filtered as a}
                    <button type="button" class="item" on:click={() => onSelect(a)}>
                        {#if isImage(a.icon)}
                            <img src={a.icon} alt={a.label} on:error={(e) => (e.target as HTMLImageElement).style.display = 'none'} />
                        {:else}
                            <span class="emoji">{a.icon}</span>
                        {/if}
                        <span class="label">{a.label}</span>
                        <span class="page-badge">{pageName(a.page)}</span>
                    </button>
                {/each}
            {/if}
        </div>

        <button class="new" on:click={() => onNew()}>＋ Nueva acción</button>
    </div>
</div>