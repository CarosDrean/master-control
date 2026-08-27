<script lang="ts">
    import { onMount } from 'svelte';
    import { createAction } from '$lib/api/actions';
    import { getActionTypes } from '$lib/api/actiontypes';
    import type { ActionType } from '$lib/api/actiontypes';
    import { getApps, appIconUrl, appActionId } from '$lib/api/apps';
    import type { App } from '$lib/api/apps';
    import { getSystemActions } from '$lib/api/system';
    import type { SystemAction } from '$lib/api/system';
    import HotkeyRecorder from '$lib/components/HotkeyRecorder.svelte';

    export let pageId = "general";
    export let position: number | null = null;
    export let onSaved: () => void = () => {};

    let actionTypes: ActionType[] = [];
    let apps: App[] = [];
    let appSearch = "";
    let appOpen = false;
    let systemActions: SystemAction[] = [];
    let systemSearch = "";
    let systemOpen = false;
    let type = "";
    let positionStr = "";
    let id = "";
    let idAuto = true;
    let label = "";
    let command = "";
    let icon = "";

    // La posición se fija al elegir un slot: no es editable ni visible.
    $: locked = position !== null && position !== undefined;

    function slugify(s: string): string {
        return s
            .toLowerCase()
            .normalize("NFD")
            .replace(/[\u0300-\u036f]/g, "")
            .replace(/[^a-z0-9]+/g, "_")
            .replace(/^_+|_+$/g, "");
    }

    function onLabelInput(e: Event) {
        const value = (e.target as HTMLInputElement).value;
        label = value;

        if (!id) {
            id = slugify(value);
            idAuto = true;
        }
    }

    function onHotkeyCaptured(info: { modifiers: string[]; key: string }) {
        if (!idAuto && id) return;

        const modNames: Record<string, string> = {
            meta: "cmd",
            control: "ctrl",
            alt: "alt",
            shift: "shift",
        };
        const modsSlug = info.modifiers.map(m => modNames[m] ?? m).join("_");
        id = ["hotkey", ...(modsSlug ? [modsSlug] : []), info.key].join("_");
        idAuto = true;
    }

    $: filteredApps = appSearch
        ? apps.filter(a => a.name.toLowerCase().includes(appSearch.toLowerCase()))
        : apps;
    $: selectedApp = apps.find(a => a.name === command);
    $: filteredSystemActions = systemSearch
        ? systemActions.filter(a => a.label.toLowerCase().includes(systemSearch.toLowerCase()))
        : systemActions;
    $: selectedSystem = systemActions.find(a => a.command === command);

    function openAppPicker() {
        appSearch = "";
        appOpen = true;
    }

    function openSystemPicker() {
        systemSearch = "";
        systemOpen = true;
    }

    function selectSystem(a: SystemAction) {
        command = a.command;
        icon = a.icon;
        label = a.label;
        id = slugify(a.label);
        idAuto = false;
        systemOpen = false;
    }

    function selectApp(a: App) {
        command = a.name;
        icon = appIconUrl(a.path);
        label = a.name;
        id = appActionId(a.name);
        idAuto = false;
        appOpen = false;
    }

    onMount(async () => {
        [actionTypes, apps, systemActions] = await Promise.all([getActionTypes(), getApps(), getSystemActions()]);
        if (actionTypes.length > 0) {
            type = actionTypes[0].id;
        }
    });

    let message = "";
    let messageType: "success" | "error" | "" = "";

    function showMessage(text: string, type: "success" | "error") {
        message = text;
        messageType = type;

        setTimeout(() => {
            message = "";
            messageType = "";
        }, 2000);
    }

    async function save() {
        if (!id) {
            id = slugify(label) || `${type}_${Date.now()}`;
        }

        try {
            const res = await createAction({
                id,
                label,
                type,
                command,
                icon,
                page: pageId,
                position: locked ? position! : (positionStr === "" ? -1 : parseInt(positionStr, 10))
            });

            if (res.success) {
                showMessage(res.message, "success");
                setTimeout(onSaved, 800);
            } else {
                showMessage(res.message, "error");
            }

        } catch (e) {
            showMessage("Error de conexión con el servidor", "error");
        }
    }

    $: isImage = icon?.startsWith("http") || icon?.startsWith("/");
    $: isAppType = type === "app";
    $: if (type === "hotkey" && !icon) {
        icon = "⌨";
    }
</script>

<style>
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
        box-sizing: border-box;
        width: 100%;
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

    .sys-icon {
        width: 26px;
        flex-shrink: 0;
        text-align: center;
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
        <label for="type">Tipo</label>
        <select id="type" bind:value={type}>
            {#each actionTypes as t}
                <option value={t.id}>{t.name} — {t.description}</option>
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
            <HotkeyRecorder bind:command={command} onCapture={onHotkeyCaptured} />
        </div>
    {:else if type === "system"}
        <div class="field">
            <label for="system-search">Acción del sistema</label>
            {#if systemOpen}
                <input id="system-search" bind:value={systemSearch} placeholder="Buscar acción…" />
                <div class="app-list">
                    {#each filteredSystemActions as sa}
                        <button
                            type="button"
                            class="app-item {command === sa.command ? 'selected' : ''}"
                            on:click={() => selectSystem(sa)}
                        >
                            <span class="sys-icon">{sa.icon}</span>
                            <span>{sa.label}</span>
                            {#if command === sa.command}
                                <span class="check">✓</span>
                            {/if}
                        </button>
                    {/each}
                    {#if filteredSystemActions.length === 0}
                        <div class="app-empty">Sin resultados</div>
                    {/if}
                </div>
            {:else}
                <button type="button" class="app-select" on:click={openSystemPicker}>
                    {#if selectedSystem}
                        <span class="sys-icon">{selectedSystem.icon}</span>
                        <span>{selectedSystem.label}</span>
                    {:else}
                        <span class="placeholder">Selecciona una acción…</span>
                    {/if}
                    <span class="chevron">▾</span>
                </button>
            {/if}
        </div>
    {:else}
        <div class="field">
            <label for="command">Command / Script</label>
            <textarea id="command" bind:value={command} placeholder="osascript -e 'tell application …'"></textarea>
        </div>
    {/if}

    <div class="field">
        <label for="label">Label</label>
        <input id="label" value={label} on:input={onLabelInput} placeholder="Nombre visible" />
    </div>

    <div class="field">
        <label for="icon">Icono (emoji o URL)</label>
        <input id="icon" bind:value={icon} placeholder="😀 o https://…" />
    </div>

    {#if !locked}
        <div class="field">
            <label for="position">Posición en la cuadrícula (vacío = sin asignar)</label>
            <input id="position" bind:value={positionStr} placeholder="0, 1, 2… (0 = primera casilla)" inputmode="numeric" />
        </div>
    {/if}

    <button class="btn-save" on:click={save}>Guardar</button>
</div>