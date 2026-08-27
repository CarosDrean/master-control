<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import ActionButton from '$lib/components/ActionButton.svelte';
    import ActionPicker from '$lib/components/ActionPicker.svelte';
    import NewActionForm from '$lib/components/NewActionForm.svelte';
    import NewPageForm from '$lib/components/NewPageForm.svelte';
    import Modal from '$lib/components/Modal.svelte';
    import { getActions, updateAction, deleteAction } from '$lib/api/actions';
    import { getPages, updatePage } from '$lib/api/pages';
    import type { Action } from '$lib/api/actions';
    import type { Page } from '$lib/api/pages';

    let actions: Action[] = [];
    let pages: Page[] = [];
    let currentPage = 0;
    let loading = true;
    let picker: { page: Page; slot: number } | null = null;
    let newAction: { page: Page; slot: number } | null = null;
    let newPage = false;
    let menu: { page: Page; slot: number; action: Action } | null = null;
    let settings: Page | null = null;
    let settingsName = "";
    let settingsIcon = "";
    let settingsColumns = 3;
    let settingsRows = 3;

    const gridCols = (p: Page) => p.columns || 3;
    const gridSlots = (p: Page) => (p.columns || 3) * (p.rows || 3);

    const actionAt = (pageId: string, slot: number) =>
        actions.find(a => (a.page || 'general') === pageId && a.position === slot);

    const isAssigned = (a: Action) => a.position !== undefined && a.position >= 0;

    // ¿Está la acción realmente en un slot válido de su página?
    const placed = (a: Action) => {
        const p = pages.find(x => x.id === (a.page || 'general'));
        return isAssigned(a) && !!p && a.position! < gridSlots(p);
    };

    async function load() {
        loading = true;
        const [pageData, actionData] = await Promise.all([getPages(), getActions()]);
        pages = pageData.sort((a, b) => a.order - b.order);
        actions = actionData;
        fillFreeSlots();
        loading = false;
    }

    onMount(load);

    // Rellena en memoria los slots libres con las acciones sin posición,
    // para que ninguna acción quede oculta. No persiste el relleno automático.
    function fillFreeSlots() {
        for (const p of pages) {
            const total = gridSlots(p);
            for (const a of actions) {
                if ((a.page || 'general') === p.id && isAssigned(a) && a.position! >= total) {
                    a.position = -1;
                }
            }

            const occupied = new Set(
                actions.filter(a => (a.page || 'general') === p.id && placed(a)).map(a => a.position)
            );
            const free: number[] = [];
            for (let i = 0; i < total; i++) {
                if (!occupied.has(i)) free.push(i);
            }

            let fi = 0;
            for (const a of actions) {
                if ((a.page || 'general') === p.id && !isAssigned(a) && fi < free.length) {
                    a.position = free[fi++];
                }
            }
        }
    }

    const actionsOnPage = (pageId: string) => actions.filter(a => (a.page || 'general') === pageId);
    const overflowOn = (pageId: string) => actionsOnPage(pageId).filter(a => !placed(a));

    function scrollTo(index: number) {
        const el = document.querySelector('.carousel') as HTMLElement | null;
        if (el) {
            const width = el.clientWidth;
            el.scrollTo({ left: index * width, behavior: 'smooth' });
        }
    }

    function onScroll() {
        const el = document.querySelector('.carousel') as HTMLElement | null;
        if (!el) return;
        currentPage = Math.round(el.scrollLeft / el.clientWidth);
    }

    function openPicker(p: Page, slot: number) {
        picker = { page: p, slot };
    }

    async function onPicked(a: Action) {
        if (!picker) return;
        a.page = picker.page.id;
        a.position = picker.slot;
        await updateAction({ ...a });
        picker = null;
    }

    async function freeAction(a: Action) {
        a.position = -1;
        await updateAction({ ...a });
    }

    async function replaceFromMenu() {
        if (!menu) return;
        const { page, slot, action } = menu;
        await freeAction(action);
        menu = null;
        openPicker(page, slot);
    }

    async function vacate() {
        if (!menu) return;
        await freeAction(menu.action);
        menu = null;
    }

    async function removeAction() {
        if (!menu) return;
        const target = menu.action;
        menu = null;
        await deleteAction(target.id);
        actions = actions.filter(a => a.id !== target.id);
    }

    function editAction() {
        if (!menu) return;
        goto(`/actions/${menu.action.id}`);
        menu = null;
    }

    function goNewAction() {
        if (!picker) return;
        const target = picker;
        picker = null;
        newAction = { page: target.page, slot: target.slot };
    }

    async function onNewSaved() {
        newAction = null;
        actions = await getActions();
        fillFreeSlots();
    }

    async function onNewPageSaved() {
        newPage = false;
        pages = (await getPages()).sort((a, b) => a.order - b.order);
        fillFreeSlots();
    }

    function openSettings(p: Page) {
        settings = p;
        settingsName = p.name;
        settingsIcon = p.icon ?? "";
        settingsColumns = p.columns || 3;
        settingsRows = p.rows || 3;
    }

    async function saveSettings() {
        if (!settings) return;
        await updatePage({ ...settings, name: settingsName, icon: settingsIcon, columns: settingsColumns, rows: settingsRows });
        pages = (await getPages()).sort((a, b) => a.order - b.order);
        settings = null;
        fillFreeSlots();
    }

    $: unassigned = actions.filter(a => !placed(a));
</script>

<style>
    .container {
        display: flex;
        flex-direction: column;
        height: 100%;
        overflow: hidden;
    }

    .carousel {
        display: flex;
        flex: 1;
        overflow-x: auto;
        scroll-snap-type: x mandatory;
        scroll-behavior: smooth;
        -webkit-overflow-scrolling: touch;
        touch-action: pan-x pan-y;
    }

    .carousel::-webkit-scrollbar {
        display: none;
    }

    .page {
        flex: 0 0 100%;
        scroll-snap-align: start;
        overflow-y: auto;
        padding: 16px 16px 8px;
        box-sizing: border-box;
    }

    .page-title {
        text-align: center;
        font-size: 1.1rem;
        opacity: 0.85;
        margin: 0 0 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
    }

    .page-settings {
        background: none;
        border: none;
        color: rgba(255,255,255,0.4);
        font-size: 1rem;
        cursor: pointer;
        padding: 0 4px;
        opacity: 0;
        transition: opacity 0.15s;
    }

    .page-title:hover .page-settings {
        opacity: 1;
    }

    .page-settings:hover {
        color: rgba(255,255,255,0.85);
    }

    .overflow {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        gap: 16px;
        max-width: 900px;
        margin: 28px auto 0;
        padding-top: 22px;
        border-top: 1px dashed rgba(255,255,255,0.12);
    }

    .overflow-title {
        width: 100%;
        text-align: center;
        font-size: 0.85rem;
        opacity: 0.5;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(var(--cols), 1fr);
        gap: 28px;
        justify-items: center;
        align-content: start;
        max-width: min(900px, calc(var(--cols) * 140px + (var(--cols) - 1) * 28px));
        margin: 0 auto;
    }

    .slot {
        width: 140px;
        height: 140px;
        border-radius: 22px;
        border: 2px dashed rgba(255,255,255,0.18);
        background: rgba(255,255,255,0.02);
        color: rgba(255,255,255,0.4);
        font-size: 40px;
        font-weight: 300;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: 0.16s ease;
        touch-action: manipulation;
    }

    .slot:hover {
        border-color: rgba(120, 180, 255, 0.45);
        color: rgba(255,255,255,0.7);
        background: rgba(120, 180, 255, 0.06);
    }

    .add-page {
        flex: 0 0 100%;
        scroll-snap-align: start;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 24px;
        box-sizing: border-box;
    }

    .add-page-btn {
        width: 220px;
        height: 220px;
        border-radius: 30px;
        border: 2px dashed rgba(255,255,255,0.25);
        background: rgba(255,255,255,0.02);
        color: rgba(255,255,255,0.5);
        font-size: 72px;
        font-weight: 200;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: 0.16s ease;
        touch-action: manipulation;
    }

    .add-page-btn:hover {
        border-color: rgba(120, 180, 255, 0.5);
        color: rgba(255,255,255,0.85);
        background: rgba(120, 180, 255, 0.06);
    }

    .empty {
        text-align: center;
        opacity: 0.5;
        margin-top: 40px;
    }

    .dots {
        display: flex;
        justify-content: center;
        gap: 10px;
        padding: 16px 0 24px;
    }

    .dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        border: none;
        padding: 0;
        background: rgba(255,255,255,0.2);
        cursor: pointer;
        transition: background 0.2s;
    }

    .dot.active {
        background: rgba(255,255,255,0.9);
    }

    .dot.add {
        background: rgba(255,255,255,0.08);
    }

    /* MODAL MENU */
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

    .menu-modal {
        width: 100%;
        max-width: 360px;
        background: #161616;
        border: 1px solid rgba(255,255,255,0.1);
        border-radius: 18px;
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
        padding: 18px;
    }

    .menu-modal h3 {
        margin: 0 0 16px;
        font-size: 1.05rem;
    }

    .menu-btn {
        display: block;
        width: 100%;
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 12px;
        padding: 12px;
        margin-bottom: 10px;
        font-size: 0.95rem;
        color: white;
        cursor: pointer;
        text-align: left;
    }

    .menu-btn:hover {
        background: #292929;
    }

    .menu-btn.danger {
        border-color: rgba(255,100,100,0.35);
        color: #ff9c9c;
    }

    .settings-row {
        display: flex;
        gap: 14px;
        margin-bottom: 16px;
    }

    .settings-row .field {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .settings-row label {
        font-size: 0.85rem;
        opacity: 0.8;
    }

    .settings-row input {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 10px 12px;
        color: white;
        font-size: 0.95rem;
        outline: none;
        width: 100%;
        box-sizing: border-box;
        min-width: 0;
    }

    .settings-row input#settings-icon {
        max-width: 96px;
    }

    .empty-btn {
        display: block;
        margin: 16px auto 0;
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.15);
        border-radius: 12px;
        padding: 10px 22px;
        font-size: 0.95rem;
        color: white;
        cursor: pointer;
    }

    .empty-btn:hover {
        background: #2a2a2a;
    }
</style>

<div class="container">
    {#if loading}
        <div class="empty">Cargando…</div>
    {:else if pages.length === 0}
        <div class="empty">
            No hay páginas todavía.
            <button class="empty-btn" on:click={() => newPage = true}>＋ Crear primera página</button>
        </div>
    {:else}
        <div class="carousel" on:scroll={onScroll}>
            {#key actions}
            {#each pages as p, i}
                {@const total = gridSlots(p)}
                {@const cols = gridCols(p)}
                <section class="page">
                    <h2 class="page-title">
                        {#if p.icon}<span>{p.icon}</span>{/if}{p.name}
                        <button class="page-settings" aria-label="Configurar página" on:click={() => openSettings(p)}>⛶</button>
                    </h2>

                    <div class="grid" style="--cols: {cols}">
                        {#each Array(total) as _, slot}
                            {@const action = actionAt(p.id, slot)}
                            {#if action}
                                <ActionButton
                                    id={action.id}
                                    label={action.label}
                                    icon={action.icon}
                                    type={action.type}
                                    command={action.command}
                                    onMenu={() => menu = { page: p, slot, action }}
                                />
                            {:else}
                                <button
                                    type="button"
                                    class="slot"
                                    aria-label="Asignar acción al slot {slot + 1}"
                                    on:click={() => openPicker(p, slot)}
                                >+</button>
                            {/if}
                        {/each}
                    </div>

                    {#if overflowOn(p.id).length > 0}
                        <div class="overflow">
                            <span class="overflow-title">Sin asignar</span>
                            {#each overflowOn(p.id) as a}
                                <ActionButton
                                    id={a.id}
                                    label={a.label}
                                    icon={a.icon}
                                    type={a.type}
                                    command={a.command}
                                    onMenu={() => menu = { page: p, slot: -1, action: a }}
                                />
                            {/each}
                        </div>
                    {/if}
                </section>
            {/each}
            {/key}

            <section class="add-page">
                <button type="button" class="add-page-btn" on:click={() => newPage = true} aria-label="Nueva página">
                    +
                </button>
            </section>
        </div>

        <div class="dots">
            {#each pages as p, i}
                <button
                    class:active={i === currentPage}
                    class="dot"
                    aria-label={p.name}
                    on:click={() => scrollTo(i)}
                ></button>
            {/each}
            <button
                class:active={pages.length === currentPage}
                class="dot add"
                aria-label="Nueva página"
                on:click={() => scrollTo(pages.length)}
            ></button>
        </div>
    {/if}
</div>

{#if picker}
    <ActionPicker
        actions={unassigned}
        pages={pages}
        onSelect={onPicked}
        onClose={() => picker = null}
        onNew={goNewAction}
    />
{/if}

{#if menu}
    <div
        class="overlay"
        role="button"
        tabindex="-1"
        aria-label="Cerrar menú"
        on:click={() => menu = null}
        on:keydown={(e) => { if (e.key === 'Escape' || e.key === 'Enter') menu = null; }}
    >
        <div class="menu-modal" role="presentation" on:click|stopPropagation>
            <h3>{menu.action.label}</h3>
            <button class="menu-btn" on:click={replaceFromMenu}>↺ Reasignar</button>
            <button class="menu-btn" on:click={editAction}>✎ Editar acción</button>
            <button class="menu-btn danger" on:click={vacate}>✕ Vaciar slot</button>
            <button class="menu-btn danger" on:click={removeAction}>🗑 Eliminar acción</button>
        </div>
    </div>
{/if}

{#if settings}
    <div
        class="overlay"
        role="button"
        tabindex="-1"
        aria-label="Cerrar configuración"
        on:click={() => settings = null}
        on:keydown={(e) => { if (e.key === 'Escape' || e.key === 'Enter') settings = null; }}
    >
        <div class="menu-modal" role="presentation" on:click|stopPropagation>
            <h3>Configurar página: {settingsName}</h3>
            <div class="settings-row">
                <div class="field">
                    <label for="settings-name">Nombre</label>
                    <input id="settings-name" bind:value={settingsName} />
                </div>
                <div class="field">
                    <label for="settings-icon">Icono (emoji)</label>
                    <input id="settings-icon" bind:value={settingsIcon} />
                </div>
            </div>
            <div class="settings-row">
                <div class="field">
                    <label for="settings-cols">Columnas</label>
                    <input id="settings-cols" type="number" min="1" max="6" bind:value={settingsColumns} />
                </div>
                <div class="field">
                    <label for="settings-rows">Filas</label>
                    <input id="settings-rows" type="number" min="1" max="6" bind:value={settingsRows} />
                </div>
            </div>
            <button class="menu-btn" on:click={saveSettings}>Guardar</button>
        </div>
    </div>
{/if}

{#if newAction}
    <Modal title="Nueva acción" onClose={() => newAction = null}>
        <NewActionForm
            pageId={newAction.page.id}
            position={newAction.slot}
            onSaved={onNewSaved}
        />
    </Modal>
{/if}

{#if newPage}
    <Modal title="Nueva página" onClose={() => newPage = false}>
        <NewPageForm onSaved={onNewPageSaved} />
    </Modal>
{/if}