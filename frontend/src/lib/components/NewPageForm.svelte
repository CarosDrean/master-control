<script lang="ts">
    import { onMount } from 'svelte';
    import { getPages, createPage } from '$lib/api/pages';

    export let onSaved: () => void = () => {};

    let id = "";
    let name = "";
    let icon = "";
    let columns = 3;
    let rows = 3;
    let order = 0;

    onMount(async () => {
        const pages = (await getPages()).sort((a, b) => a.order - b.order);
        order = pages.length;
    });

    async function save() {
        if (!id || !name) return;
        await createPage({ id, name, icon, order, columns, rows });
        onSaved();
    }
</script>

<style>
    .field {
        display: flex;
        flex-direction: column;
        margin-bottom: 14px;
    }

    .field-row {
        display: flex;
        gap: 14px;
    }

    .field-row .field {
        flex: 1;
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
        box-sizing: border-box;
        width: 100%;
    }

    input:focus {
        border-color: rgba(255,255,255,0.25);
    }

    .btn {
        width: 100%;
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 12px;
        padding: 12px;
        font-size: 1rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
        margin-top: 12px;
    }

    .btn:hover {
        background: #2a2a2a;
    }

    .btn:active {
        transform: scale(0.97);
    }
</style>

<div class="field">
    <label for="page-id">ID</label>
    <input id="page-id" bind:value={id} placeholder="general" />
</div>
<div class="field">
    <label for="page-name">Nombre</label>
    <input id="page-name" bind:value={name} placeholder="General" />
</div>
<div class="field">
    <label for="page-icon">Icono (emoji)</label>
    <input id="page-icon" bind:value={icon} placeholder="✨" />
</div>
<div class="field-row">
    <div class="field">
        <label for="page-cols">Columnas</label>
        <input id="page-cols" type="number" min="1" max="6" bind:value={columns} />
    </div>
    <div class="field">
        <label for="page-rows">Filas</label>
        <input id="page-rows" type="number" min="1" max="6" bind:value={rows} />
    </div>
</div>

<button class="btn" on:click={save}>Añadir página</button>