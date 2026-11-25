<script>
    import '../app.css';
    import { goto } from '$app/navigation';
    import SystemHud from '$lib/components/SystemHud.svelte';
    import ActionFeedback from '$lib/components/ActionFeedback.svelte';
    import TokenGate from '$lib/components/TokenGate.svelte';
    import { onMount } from 'svelte';
    import { keepAwake } from '$lib/keepAwake';

    // La pantalla no se apaga mientras Master Control está abierto.
    onMount(keepAwake);
</script>

<style>
    header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        padding: 16px 24px;
        background: #151515;
        border-bottom: 1px solid rgba(255, 255, 255, 0.06);

        position: sticky;
        top: 0;
        z-index: 5;
        color: white;
    }

    .left {
        display: flex;
        flex-direction: column;
    }

    .title-btn {
        background: none;
        border: none;
        color: white;
        font-size: 1.4rem;
        font-weight: 600;
        cursor: pointer;
        padding: 0;
    }

    .nav {
        color: white;
        text-decoration: none;
        font-size: 1rem;
        padding: 8px 14px;
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.12);
        background: #1f1f1f;
    }

    main {
        padding: 0;
        height: calc(100dvh - 61px);
        overflow-y: auto;
    }
</style>

<header>
    <div class="left">
        <button class="title-btn" type="button" on:click={() => goto('/')}>
            Master Control
        </button>
    </div>
    <a class="nav" href="/scroll">↕ Scroll</a>
</header>

<main>
    <TokenGate>
        <slot />
    </TokenGate>
</main>

<SystemHud />
<ActionFeedback />