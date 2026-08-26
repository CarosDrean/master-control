<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { getAction, updateAction, executeAction } from '$lib/api/actions';
    import { getActionTypes } from '$lib/api/actiontypes';
    import type { ActionType } from '$lib/api/actiontypes';
    import { getPages } from '$lib/api/pages';
    import type { Page } from '$lib/api/pages';
    import { getApps, appIconUrl, appActionId } from '$lib/api/apps';
    import type { App } from '$lib/api/apps';
    import HotkeyRecorder from '$lib/components/HotkeyRecorder.svelte';

    export let params;

    let actionTypes: ActionType[] = [];
    let pages: Page[] = [];
    let apps: App[] = [];
    let appSearch = "";
    let appOpen = false;

    let action: any = null;

    let id = "";
    let label = "";
    let type = "";
    let command = "";
    let icon = "";
    let page = "general";

    let message = "";
    let messageType: "success" | "error" | "" = "";

    $: filteredApps = appSearch
        ? apps.filter(a => a.name.toLowerCase().includes(appSearch.toLowerCase()))
        : apps;
    $: selectedApp = apps.find(a => a.name === command);

    function openAppPicker() {
        appSearch = "";
        appOpen = true;
    }

    function selectApp(a: App) {
        command = a.name;
        icon = appIconUrl(a.path);
        label = a.name;
        id = appActionId(a.name);
        appOpen = false;
    }

    function showMessage(text: string, type: "success" | "error") {
        message = text;
        messageType = type;

        setTimeout(() => {
            message = "";
            messageType = "";
        }, 2000);
    }

    onMount(async () => {
        [action, actionTypes, pages, apps] = await Promise.all([getAction(params.id), getActionTypes(), getPages(), getApps()]);
        pages = pages.sort((a, b) => a.order - b.order);

        id = action.id;
        label = action.label;
        type = action.type;
        command = action.command;
        icon = action.icon ?? "";
        page = action.page ?? "general";
    });

    async function save() {
        try {
            const res = await updateAction({
                id,
                label,
                type,
                command,
                icon,
                page
            });

            if (res.success) {
                showMessage(res.message, "success");
                setTimeout(() => goto('/actions'), 1200);
            } else {
                showMessage(res.message, "error");
            }

        } catch (e) {
            showMessage("Error de conexión con el servidor", "error");
        }
    }

    async function test() {
        try {
            await executeAction(id);
            showMessage("Acción ejecutada", "success");
        } catch (e) {
            showMessage("Error al ejecutar", "error");
        }
    }

    $: isImage = icon?.startsWith("http") || icon?.startsWith("/");
    $: isAppType = type === "app";
</script>

<style>
    .container {
        max-width: 600px;
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
        padding: 30px;
    }

    .field {
        display: flex;
        flex-direction: column;
        margin-bottom: 16px;
    }

    label {
        font-size: 0.85rem;
        margin-bottom: 6px;
        opacity: 0.8;
    }

    input, select, textarea {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 12px 14px;
        color: white;
        font-size: 0.95rem;
        outline: none;
        transition: border 0.15s;
    }

    input:focus,
    select:focus,
    textarea:focus {
        border-color: rgba(255,255,255,0.25);
    }

    textarea {
        min-height: 120px;
        resize: vertical;
    }

    .app-list {
        margin-top: 8px;
        max-height: 240px;
        overflow-y: auto;
        border: 1px solid rgba(255,255,255,0.1);
        border-radius: 10px;
        background: #1a1a1a;
        padding: 4px;
    }

    .app-select {
        display: flex;
        align-items: center;
        gap: 10px;
        width: 100%;
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 8px 14px;
        color: white;
        font-size: 0.95rem;
        cursor: pointer;
        transition: border 0.15s;
    }

    .app-select:hover {
        border-color: rgba(255,255,255,0.25);
    }

    .app-select img {
        width: 26px;
        height: 26px;
        border-radius: 6px;
        flex-shrink: 0;
    }

    .app-select .placeholder {
        opacity: 0.5;
    }

    .app-select .chevron {
        margin-left: auto;
        opacity: 0.6;
    }

    .app-item {
        display: flex;
        align-items: center;
        gap: 10px;
        width: 100%;
        background: none;
        border: none;
        border-radius: 8px;
        padding: 8px 10px;
        color: white;
        font-size: 0.9rem;
        text-align: left;
        cursor: pointer;
    }

    .app-item:hover {
        background: rgba(255,255,255,0.06);
    }

    .app-item.selected {
        background: rgba(255,255,255,0.12);
    }

    .app-item img {
        width: 26px;
        height: 26px;
        border-radius: 6px;
        flex-shrink: 0;
    }

    .app-item .check {
        margin-left: auto;
        opacity: 0.8;
    }

    .app-empty {
        padding: 14px;
        text-align: center;
        opacity: 0.5;
        font-size: 0.9rem;
    }

    /* PREVIEW STYLE */
    .preview {
        width: 140px;
        height: 140px;
        border-radius: 20px;
        background: linear-gradient(180deg, #1c1c1c, #111);
        border: 1px solid rgba(255,255,255,0.1);
        box-shadow:
                0 6px 16px rgba(0,0,0,0.5),
                inset 0 0 0 3px rgba(255,255,255,0.05);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin: 20px auto;
    }

    .preview img {
        width: 48px;
        height: 48px;
    }

    .preview .emoji {
        font-size: 52px;
    }

    .preview-label {
        margin-top: 6px;
        font-size: 0.9rem;
        opacity: 0.9;
        text-align: center;
    }

    .btn-save {
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

    .btn-save:hover {
        background: #292929;
    }

    .btn-save:active {
        transform: scale(0.97);
    }

    .btn-back {
        background: none;
        border: none;
        color: #aaa;
        cursor: pointer;
        font-size: 0.9rem;
        margin-bottom: 20px;
    }

    .toast {
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 14px 20px;
        border-radius: 10px;
        color: white;
        font-size: 0.95rem;
        background: rgba(30, 30, 30, 0.9);
        border: 1px solid rgba(255,255,255,0.12);
        box-shadow: 0 4px 12px rgba(0,0,0,0.4);
        animation: fadeIn 0.3s ease;
        z-index: 999;
    }

    .toast.success {
        border-color: rgba(0,255,150,0.5);
    }

    .toast.error {
        border-color: rgba(255,100,100,0.5);
    }

    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(-10px); }
        to   { opacity: 1; transform: translateY(0); }
    }

</style>

<div class="container">
    <button class="btn-back" on:click={() => goto('/actions')}>
        ← Regresar
    </button>

    <h1>Editar acción</h1>

    <!-- PREVIEW -->
    <div class="preview">
        {#if isImage}
            <img src={icon} alt="" />
        {:else}
            <div class="emoji">{icon}</div>
        {/if}

        <div class="preview-label">{label}</div>
    </div>

    {#if message}
        <div class="toast {messageType}">
            {message}
        </div>
    {/if}

    <div class="panel">
        <div class="field">
            <label for="label-edit">Label</label>
            <input id="label-edit" bind:value={label} placeholder="Nombre visible" />
        </div>

        <div class="field">
            <label for="icon-edit">Icono (emoji o URL)</label>
            <input id="icon-edit" bind:value={icon} placeholder="😀 o https://…" />
        </div>

        <div class="field">
            <label for="type-edit">Tipo</label>
            <select id="type-edit" bind:value={type}>
                {#each actionTypes as t}
                    <option value={t.id}>{t.name} — {t.description}</option>
                {/each}
            </select>
        </div>

        <div class="field">
            <label for="page-edit">Página</label>
            <select id="page-edit" bind:value={page}>
                {#each pages as p}
                    <option value={p.id}>{p.icon} {p.name}</option>
                {/each}
            </select>
        </div>

        {#if isAppType}
            <div class="field">
                <label for="app-search">App</label>
                {#if appOpen}
                    <input id="app-search" bind:value={appSearch} placeholder="Buscar app…" />
                    <div class="app-list">
                        {#each filteredApps as a}
                            <button
                                type="button"
                                class="app-item {command === a.name ? 'selected' : ''}"
                                on:click={() => selectApp(a)}
                            >
                                <img src={appIconUrl(a.path)} alt={a.name} on:error={(e) => (e.target as HTMLImageElement).style.display = 'none'} />
                                <span>{a.name}</span>
                                {#if command === a.name}
                                    <span class="check">✓</span>
                                {/if}
                            </button>
                        {/each}
                        {#if filteredApps.length === 0}
                            <div class="app-empty">Sin resultados</div>
                        {/if}
                    </div>
                {:else}
                    <button type="button" class="app-select" on:click={openAppPicker}>
                        {#if selectedApp}
                            <img src={appIconUrl(selectedApp.path)} alt={selectedApp.name} on:error={(e) => (e.target as HTMLImageElement).style.display = 'none'} />
                            <span>{selectedApp.name}</span>
                        {:else}
                            <span class="placeholder">Selecciona una app…</span>
                        {/if}
                        <span class="chevron">▾</span>
                    </button>
                {/if}
            </div>
        {:else if type === "hotkey"}
            <div class="field">
                <label for="hotkey-command">Hotkey</label>
                <HotkeyRecorder bind:command={command} />
            </div>
        {:else}
            <div class="field">
                <label for="command-edit">Command / Script</label>
                <textarea id="command-edit" bind:value={command}></textarea>
            </div>
        {/if}

        <button class="btn-save" on:click={test}>▶ Probar</button>
        <button class="btn-save" on:click={save}>Guardar cambios</button>
    </div>
</div>
