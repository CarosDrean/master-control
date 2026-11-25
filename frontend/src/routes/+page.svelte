<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import ActionButton from '$lib/components/ActionButton.svelte';
    import ActionPicker from '$lib/components/ActionPicker.svelte';
    import NewActionForm from '$lib/components/NewActionForm.svelte';
    import NewPageForm from '$lib/components/NewPageForm.svelte';
    import Modal from '$lib/components/Modal.svelte';
    import { getActions, updateAction, deleteAction } from '$lib/api/actions';
    import { getPages, updatePage, createPage } from '$lib/api/pages';
    import { getApps, type App } from '$lib/api/apps';
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
    let settingsApp = "";
    let installedApps: App[] = [];

    // ── Cajón de una app ──────────────────────────────────────────────────
    // Tocar el botón de una app la abre en la Mac y despliega, al costado de la cuadrícula, las acciones
    // de esa app: las de la página vinculada a ella (Page.app). Esas páginas no salen en el
    // carrusel: son el "detalle" de su app.
    let expanded: { page: string; slot: number; app: string; label: string } | null = null;

    $: visiblePages = pages.filter(p => !p.app);
    // Una app puede tener varias páginas de acciones: el cajón las pasa deslizando.
    const appPages = (app: string) => pages.filter(p => p.app === app).sort((a, b) => a.order - b.order);

    let drawerEl: HTMLElement | null = null;
    let drawerIndex = 0;

    function toggleDrawer(p: Page, slot: number, a: Action) {
        const same = expanded && expanded.page === p.id && expanded.slot === slot;
        expanded = same ? null : { page: p.id, slot, app: a.command, label: a.label };
        drawerIndex = 0;
    }

    function onDrawerScroll() {
        if (drawerEl) drawerIndex = Math.round(drawerEl.scrollLeft / drawerEl.clientWidth);
    }

    function drawerTo(index: number) {
        drawerEl?.scrollTo({ left: index * drawerEl.clientWidth, behavior: 'smooth' });
    }

    async function createAppPage() {
        if (!expanded) return;
        const { app, label } = expanded;
        const slug = app.toLowerCase().replace(/[^a-z0-9]+/g, "_").replace(/^_+|_+$/g, "");
        const existing = appPages(app);
        const n = existing.length + 1;
        let id = n === 1 ? `app_${slug}` : `app_${slug}_${n}`;
        while (pages.some(x => x.id === id)) id += "_";
        const order = Math.max(0, ...pages.map(x => x.order)) + 1;
        const cols = existing[0]?.columns || 3;
        const rows = existing[0]?.rows || 3;
        await createPage({ id, name: n === 1 ? label : `${label} ${n}`, icon: "", order, columns: cols, rows, app });
        pages = (await getPages()).sort((a, b) => a.order - b.order);
        fillFreeSlots();
        // Mostrar la página recién creada.
        setTimeout(() => drawerTo(existing.length), 0);
    }


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
        settingsApp = p.app ?? "";
        if (installedApps.length === 0) {
            getApps().then(list => installedApps = list.sort((a, b) => a.name.localeCompare(b.name))).catch(() => {});
        }
    }

    async function saveSettings() {
        if (!settings) return;
        await updatePage({ ...settings, name: settingsName, icon: settingsIcon, columns: settingsColumns, rows: settingsRows, app: settingsApp || undefined });
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

    /* El título de la página y la cabecera del cajón miden lo mismo (--head-h + margen),
       así las filas de botones de ambos lados quedan a la misma altura. */
    .page-title {
        text-align: center;
        font-size: 1.1rem;
        opacity: 0.85;
        height: var(--head-h);
        margin: 0 0 var(--head-gap);
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

    /* Cuadrícula y cajón lado a lado; si no caben (iPad en vertical), el cajón baja. */
    .page-body {
        --head-h: 32px;
        --head-gap: 16px;
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        align-items: flex-start;
        gap: var(--tile-gap);
    }

    .page-main > .grid {
        margin: 0;
    }

    /* Alineado con la página: el cajón sube lo que mide su padding, así su cabecera queda a la
       altura del título y sus filas de botones a la de las filas de apps (mismo gap). */
    .drawer {
        margin-top: -16px;
        padding: 16px;
        border-radius: 24px;
        background: rgba(120, 180, 255, 0.05);
        border: 1px solid rgba(120, 180, 255, 0.25);
        animation: drawer-in 0.18s ease-out;
    }

    @keyframes drawer-in {
        from { opacity: 0; transform: translateX(-12px); }
        to { opacity: 1; transform: none; }
    }

    .drawer .grid {
        gap: var(--tile-gap);
    }

    /* Páginas del cajón: deslizables con snap, del ancho de una cuadrícula. */
    .drawer-carousel {
        display: flex;
        width: calc(var(--dcols) * var(--tile) + (var(--dcols) - 1) * var(--tile-gap));
        max-width: calc(100vw - 64px);
        overflow-x: auto;
        scroll-snap-type: x mandatory;
        scrollbar-width: none;
        -webkit-overflow-scrolling: touch;
        overscroll-behavior-x: contain;
    }

    .drawer-carousel::-webkit-scrollbar {
        display: none;
    }

    .drawer-page {
        flex: 0 0 100%;
        scroll-snap-align: start;
        /* Que la sombra de los botones de una página no asome en el borde de la vecina;
           el padding/margin de abajo deja ver la sombra inferior sin mover el layout. */
        overflow: hidden;
        padding-bottom: 16px;
        margin-bottom: -16px;
    }

    .drawer-page .grid {
        margin: 0;
    }

    .drawer-add {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 12px;
        min-height: var(--tile);
        font-size: 0.85rem;
        color: rgba(255,255,255,0.45);
    }

    .drawer-dots {
        display: flex;
        gap: 6px;
        align-items: center;
        justify-content: center;
        margin-top: 14px;
    }

    .drawer-head {
        display: flex;
        align-items: center;
        gap: 8px;
        height: var(--head-h);
        margin-bottom: var(--head-gap);
        font-size: 1.1rem;
        opacity: 0.9;
    }

    .drawer-head span {
        flex: 1;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .drawer-btn {
        width: var(--head-h);
        height: var(--head-h);
        border-radius: 12px;
        border: 1px solid rgba(255,255,255,0.15);
        background: #1a1a1a;
        color: rgba(255,255,255,0.8);
        font-size: 1rem;
        cursor: pointer;
        touch-action: manipulation;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(var(--cols), 1fr);
        gap: var(--tile-gap);
        justify-items: center;
        align-content: start;
        max-width: min(900px, calc(var(--cols) * var(--tile) + (var(--cols) - 1) * var(--tile-gap)));
        margin: 0 auto;
    }

    .slot {
        width: var(--tile);
        height: var(--tile);
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

    .menu-modal > .field {
        display: flex;
        flex-direction: column;
        gap: 6px;
        margin-bottom: 14px;
    }

    .menu-modal > .field label {
        font-size: 0.85rem;
        opacity: 0.8;
    }

    .menu-modal > .field select {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 10px 12px;
        color: white;
        font-size: 0.95rem;
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
            {#each visiblePages as p, i}
                {@const total = gridSlots(p)}
                {@const cols = gridCols(p)}
                <section class="page">
                    <div class="page-body">
                    <div class="page-main">
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
                                    confirm={!!action.confirm}
                                    showOutput={!!action.show_output}
                                    onMenu={() => menu = { page: p, slot, action }}
                                    expanded={!!expanded && expanded.page === p.id && expanded.slot === slot}
                                    onRun={action.type === 'app' ? () => toggleDrawer(p, slot, action) : () => {}}
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
                    </div>
                    {#if expanded && expanded.page === p.id}
                        {@const aps = appPages(expanded.app)}
                        {@const dcols = Math.max(3, ...aps.map(gridCols))}
                        <div class="drawer">
                            <div class="drawer-head">
                                <span>{expanded.label}{#if drawerIndex < aps.length && aps.length > 1} · {aps[drawerIndex].name}{/if}</span>
                                {#if drawerIndex < aps.length}
                                    <button class="drawer-btn" aria-label="Configurar esta página de {expanded.label}" on:click={() => openSettings(aps[drawerIndex])}>⛶</button>
                                {/if}
                                <button class="drawer-btn" aria-label="Cerrar" on:click={() => expanded = null}>✕</button>
                            </div>
                            {#if aps.length > 0}
                                <div class="drawer-carousel" style="--dcols: {dcols}" bind:this={drawerEl} on:scroll={onDrawerScroll}>
                                    {#each aps as ap}
                                        <div class="drawer-page">
                                            <div class="grid" style="--cols: {gridCols(ap)}">
                                                {#each Array(gridSlots(ap)) as _, aslot}
                                                    {@const sub = actionAt(ap.id, aslot)}
                                                    {#if sub}
                                                        <ActionButton
                                                            id={sub.id}
                                                            label={sub.label}
                                                            icon={sub.icon}
                                                            type={sub.type}
                                                            command={sub.command}
                                                            confirm={!!sub.confirm}
                                                            showOutput={!!sub.show_output}
                                                            onMenu={() => menu = { page: ap, slot: aslot, action: sub }}
                                                        />
                                                    {:else}
                                                        <button
                                                            type="button"
                                                            class="slot"
                                                            aria-label="Asignar acción al slot {aslot + 1}"
                                                            on:click={() => openPicker(ap, aslot)}
                                                        >+</button>
                                                    {/if}
                                                {/each}
                                            </div>
                                        </div>
                                    {/each}
                                    <div class="drawer-page drawer-add">
                                        <button type="button" class="slot" aria-label="Nueva página de acciones para {expanded.label}" on:click={createAppPage}>+</button>
                                        <span>Otra página para {expanded.label}</span>
                                    </div>
                                </div>
                            {:else}
                                <button class="empty-btn" on:click={createAppPage}>＋ Agregar acciones para {expanded.label}</button>
                            {/if}
                            {#if aps.length > 0}
                                <div class="drawer-dots">
                                    {#each aps as ap, i}
                                        <button class="dot" class:active={i === drawerIndex} aria-label={ap.name} on:click={() => drawerTo(i)}></button>
                                    {/each}
                                    <button class="dot add" class:active={drawerIndex === aps.length} aria-label="Nueva página" on:click={() => drawerTo(aps.length)}></button>
                                </div>
                            {/if}
                        </div>
                    {/if}
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
                                    confirm={!!a.confirm}
                                    showOutput={!!a.show_output}
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
            {#each visiblePages as p, i}
                <button
                    class:active={i === currentPage}
                    class="dot"
                    aria-label={p.name}
                    on:click={() => scrollTo(i)}
                ></button>
            {/each}
            <button
                class:active={visiblePages.length === currentPage}
                class="dot add"
                aria-label="Nueva página"
                on:click={() => scrollTo(visiblePages.length)}
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
            <div class="field">
                <label for="settings-app">App vinculada (la página se vuelve el cajón de esa app y sale del carrusel)</label>
                <select id="settings-app" bind:value={settingsApp}>
                    <option value="">— Ninguna —</option>
                    {#if settingsApp && !installedApps.some(a => a.name === settingsApp)}
                        <option value={settingsApp}>{settingsApp}</option>
                    {/if}
                    {#each installedApps as a}
                        <option value={a.name}>{a.name}</option>
                    {/each}
                </select>
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