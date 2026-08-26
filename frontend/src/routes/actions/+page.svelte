<script lang="ts">
    import { onMount } from 'svelte';
    import { getActions, deleteAction } from '$lib/api/actions';
    import type { Action } from '$lib/api/actions';
    import { goto } from '$app/navigation';

    let actions: Action[] = [];

    onMount(async () => {
        actions = await getActions();
    });

    async function remove(id: string) {
        await deleteAction(id);
        actions = actions.filter(a => a.id !== id);
    }

    const isImage = (icon?: string) =>
        !!icon && (icon.startsWith("http") || icon.startsWith("/") || /\.(png|jpe?g|svg|webp|gif|ico)$/i.test(icon));
</script>

<style>
    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 40px 24px;
    }

    .title {
        font-size: 2rem;
        font-weight: 700;
        line-height: 1.1;
        margin-bottom: 0;
    }

    .header-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 32px;
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

    /* TABLE PANEL */
    .panel {
        margin-top: 40px;
        background: #141414;
        border-radius: 14px;
        border: 1px solid rgba(255,255,255,0.08);
        padding: 24px;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th {
        text-align: left;
        padding: 12px 8px;
        font-size: 0.85rem;
        opacity: 0.7;
        border-bottom: 1px solid rgba(255,255,255,0.08);
    }

    td {
        padding: 14px 8px;
        border-bottom: 1px solid rgba(255,255,255,0.06);
        font-size: 0.95rem;
    }

    tr:hover td {
        background: rgba(255,255,255,0.03);
    }

    .row-icon {
        width: 28px;
        height: 28px;
        object-fit: contain;
        vertical-align: middle;
    }

    .row-emoji {
        font-size: 24px;
        line-height: 1;
    }

    .btn-small {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        padding: 6px 16px;
        border-radius: 12px;
        font-size: 0.82rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
        box-shadow:
                inset 0 0 0 1px rgba(255,255,255,0.06),
                0 4px 8px rgba(0,0,0,0.35);
    }

    .btn-small:hover {
        background: #2a2a2a;
    }

    .btn-small:active {
        transform: scale(0.96);
        box-shadow:
                inset 0 0 0 1px rgba(255,255,255,0.03),
                0 2px 6px rgba(0,0,0,0.25);
    }
</style>

<div class="container">
    <div class="header-row">
        <div class="title">Administrar acciones</div>

        <button class="btn-create" on:click={() => goto('/actions/new')}>
            Nueva acción
        </button>
    </div>

    <div class="panel">
        <table>
            <thead>
            <tr>
                <th>ID</th>
                <th>Icono</th>
                <th>Label</th>
                <th>Tipo</th>
                <th>Página</th>
                <th style="text-align:right">Acciones</th>
            </tr>
            </thead>

            <tbody>
            {#each actions as a}
                <tr>
                    <td>{a.id}</td>
                    <td>
                        {#if isImage(a.icon)}
                            <img class="row-icon" src={a.icon} alt={a.label} />
                        {:else}
                            <span class="row-emoji">{a.icon}</span>
                        {/if}
                    </td>
                    <td>{a.label}</td>
                    <td>{a.type}</td>
                    <td>{a.page ?? 'general'}</td>
                    <td style="text-align:right">
                        <button class="btn-small" on:click={() => goto(`/actions/${a.id}`)}>Editar</button>
                        <button class="btn-small" on:click={() => remove(a.id)}>Eliminar</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>
