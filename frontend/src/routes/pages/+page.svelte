<script lang="ts">
    import { onMount } from 'svelte';
    import { getPages, createPage, updatePage, deletePage } from '$lib/api/pages';
    import type { Page } from '$lib/api/pages';
    import { goto } from '$app/navigation';

    let pages: Page[] = [];
    let id = "";
    let name = "";
    let icon = "";
    let editing: string | null = null;
    let editName = "";
    let editIcon = "";

    onMount(async () => {
        pages = (await getPages()).sort((a, b) => a.order - b.order);
    });

    async function add() {
        if (!id || !name) return;
        await createPage({ id, name, icon, order: pages.length });
        pages = (await getPages()).sort((a, b) => a.order - b.order);
        id = ""; name = ""; icon = "";
    }

    function startEdit(p: Page) {
        editing = p.id;
        editName = p.name;
        editIcon = p.icon ?? "";
    }

    async function saveEdit(p: Page) {
        await updatePage({ ...p, name: editName, icon: editIcon });
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
        margin-bottom: 24px;
    }

    .panel {
        background: #141414;
        border-radius: 18px;
        border: 1px solid rgba(255,255,255,0.06);
        padding: 24px;
        margin-bottom: 24px;
    }

    .field {
        display: flex;
        flex-direction: column;
        margin-bottom: 14px;
    }

    label {
        font-size: 0.85rem;
        margin-bottom: 6px;
        opacity: 0.8;
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
    <h1>Páginas</h1>

    <div class="panel">
        <div class="field">
            <label for="id">ID</label>
            <input id="id" bind:value={id} placeholder="general" />
        </div>
        <div class="field">
            <label for="name">Nombre</label>
            <input id="name" bind:value={name} placeholder="General" />
        </div>
        <div class="field">
            <label for="icon">Icono (emoji)</label>
            <input id="icon" bind:value={icon} placeholder="✨" />
        </div>
        <button class="btn" on:click={add}>Añadir página</button>
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
                        <button class="btn" on:click={() => saveEdit(p)}>Guardar</button>
                    {:else}
                        <span class="page-name">{p.name}</span>
                        <button class="btn" on:click={() => startEdit(p)}>Editar</button>
                        <button class="btn" on:click={() => remove(p)}>Eliminar</button>
                    {/if}
                </div>
            {/each}
        {/if}
    </div>
</div>