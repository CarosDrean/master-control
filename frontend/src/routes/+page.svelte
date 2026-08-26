<script lang="ts">
    import { onMount } from 'svelte';
    import ActionButton from '$lib/components/ActionButton.svelte';
    import { getActions } from '$lib/api/actions';
    import { getPages } from '$lib/api/pages';
    import type { Action } from '$lib/api/actions';
    import type { Page } from '$lib/api/pages';

    let actions: Action[] = [];
    let pages: Page[] = [];
    let currentPage = 0;
    let loading = true;

    const byPage = (pageId: string) => actions.filter(a => (a.page || 'general') === pageId);

    async function load() {
        loading = true;
        const [pageData, actionData] = await Promise.all([getPages(), getActions()]);
        pages = pageData.sort((a, b) => a.order - b.order);
        actions = actionData;
        loading = false;
    }

    onMount(load);

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
        padding: 24px 16px 8px;
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

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
        gap: 28px;
        justify-items: center;
        max-width: 900px;
        margin: 0 auto;
    }

    .empty {
        text-align: center;
        opacity: 0.5;
        margin-top: 60px;
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
</style>

<div class="container">
    {#if loading}
        <div class="empty">Cargando…</div>
    {:else if pages.length === 0}
        <div class="empty">No hay páginas. Crea una en administración.</div>
    {:else}
        <div class="carousel" on:scroll={onScroll}>
            {#each pages as p, i}
                <section class="page">
                    <h2 class="page-title">
                        {#if p.icon}<span>{p.icon}</span>{/if}{p.name}
                    </h2>

                    {#if byPage(p.id).length === 0}
                        <div class="empty">Sin acciones en esta página</div>
                    {:else}
                        <div class="grid">
                            {#each byPage(p.id) as a}
                                <ActionButton id={a.id} label={a.label} icon={a.icon} type={a.type} command={a.command} />
                            {/each}
                        </div>
                    {/if}
                </section>
            {/each}
        </div>

        {#if pages.length > 1}
            <div class="dots">
                {#each pages as p, i}
                    <button
                        class:active={i === currentPage}
                        class="dot"
                        aria-label={p.name}
                        on:click={() => scrollTo(i)}
                    ></button>
                {/each}
            </div>
        {/if}
    {/if}
</div>