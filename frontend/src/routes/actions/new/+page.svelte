<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { page as pageStore } from '$app/stores';
    import NewActionForm from '$lib/components/NewActionForm.svelte';

    let pageId = "general";
    let position: number | null = null;
    let loaded = false;

    onMount(() => {
        const params = $pageStore.url.searchParams;
        const pageParam = params.get('page');
        if (pageParam) pageId = pageParam;
        const positionParam = params.get('position');
        if (positionParam !== null) position = parseInt(positionParam, 10);
        loaded = true;
    });
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
    <button class="btn-back" on:click={() => goto('/actions')}>
        ← Regresar
    </button>

    <h1>Nueva acción</h1>

    {#if loaded}
        <NewActionForm
            pageId={pageId}
            position={position}
            onSaved={() => goto('/actions')}
        />
    {/if}
</div>