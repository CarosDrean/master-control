<script lang="ts">
    import { onMount } from 'svelte';
    import { getPages, updatePage, deletePage } from '$lib/api/pages';
    import type { Page } from '$lib/api/pages';
    import { goto } from '$app/navigation';
    import Modal from '$lib/components/Modal.svelte';
    import NewPageForm from '$lib/components/NewPageForm.svelte';

    let pages: Page[] = [];
    let showNewPage = false;
    let editing: string | null = null;
    let editName = "";
    let editIcon = "";
    let editColumns = 3;
    let editRows = 3;

    onMount(async () => {
        pages = (await getPages()).sort((a, b) => a.order - b.order);
    });

    async function onNewPageSaved() {
        showNewPage = false;
        pages = (await getPages()).sort((a, b) => a.order - b.order);
    }

    function startEdit(p: Page) {
        editing = p.id;
        editName = p.name;
        editIcon = p.icon ?? "";
        editColumns = p.columns || 3;
        editRows = p.rows || 3;
    }

    async function saveEdit(p: Page) {
        await updatePage({ ...p, name: editName, icon: editIcon, columns: editColumns, rows: editRows });
        editing = null;
        pages = (await getPages()).sort((a, b) => a.order - b.order);
    }

    async function remove(p: Page) {
        await deletePage(p.id);
        pages = pages.filter(x => x.id !== p.id);
    }
</script>

<style>
    .container {
        max-width: 700px;
        margin: 0 auto;
        padding: 40px 24px;
    }

    h1 {
        font-size: 2rem;
        font-weight: 700;
        margin-bottom: 0;
    }

    .header-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 24px;
        gap: 20px;
    }

    .btn-create {
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        padding: 10px 22px;
        border-radius: 14px;
        font-size: 0.95rem;
        color: white;
        cursor: pointer;
        transition: 0.2s;
        height: 42px;
        display: flex;
        align-items: center;
    }

    .btn-create:hover {
        background: #2a2a2a;
    }

    .panel {
        background: #141414;
        border-radius: 18px;
        border: 1px solid rgba(255,255,255,0.06);
        padding: 24px;
        margin-bottom: 24px;
    }

    input {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 12px 14px;
        color: white;
        font-size: 0.95rem;
        outline: none;
    }

    input:focus {
        border-color: rgba(255,255,255,0.25);
    }

    .btn {
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 12px;
        padding: 10px 18px;
        font-size: 0.95rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
    }

    .btn:hover {
        background: #2a2a2a;
    }

    .row {
        display: flex;
        align-items: center;
        gap: 14px;
        padding: 12px 0;
        border-bottom: 1px solid rgba(255,255,255,0.06);
    }

    .row:last-child {
        border-bottom: none;
    }

    .row .page-icon {
        font-size: 1.4rem;
        width: 32px;
        text-align: center;
    }

    .row .page-name {
        flex: 1;
    }

    .btn-back {
        background: none;
        border: none;
        color: #aaa;
        cursor: pointer;
        font-size: 0.9rem;
        margin-bottom: 20px;
    }
</style>

<div class="container">
    <button class="btn-back" on:click={() => goto('/')}>← Volver</button>

    <div class="header-row">
        <h1>Páginas</h1>
        <button class="btn-create" on:click={() => showNewPage = true}>Añadir página</button>
    </div>

    <div class="panel">
        {#if pages.length === 0}
            <div>No hay páginas</div>
        {:else}
            {#each pages as p}
                <div class="row">
                    <span class="page-icon">{p.icon}</span>
                    {#if editing === p.id}
                        <input bind:value={editName} placeholder="Nombre" />
                        <input bind:value={editIcon} placeholder="Icono" style="width:80px" />
                        <input type="number" min="1" max="6" bind:value={editColumns} style="width:70px" aria-label="Columnas" />
                        <span style="opacity:0.5">×</span>
                        <input type="number" min="1" max="6" bind:value={editRows} style="width:70px" aria-label="Filas" />
                        <button class="btn" on:click={() => saveEdit(p)}>Guardar</button>
                    {:else}
                        <span class="page-name">{p.name}</span>
                        <span style="opacity:0.5;font-size:0.85rem">{p.columns || 3}×{p.rows || 3}</span>
                        <button class="btn" on:click={() => startEdit(p)}>Editar</button>
                        <button class="btn" on:click={() => remove(p)}>Eliminar</button>
                    {/if}
                </div>
            {/each}
        {/if}
    </div>
</div>

{#if showNewPage}
    <Modal title="Nueva página" onClose={() => showNewPage = false}>
        <NewPageForm onSaved={onNewPageSaved} />
    </Modal>
{/if}