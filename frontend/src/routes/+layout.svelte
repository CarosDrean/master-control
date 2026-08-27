<script>
    import '../app.css';
    import { goto } from '$app/navigation';
    import SystemHud from '$lib/components/SystemHud.svelte';
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
</header>

<main>
    <slot />
</main>

<SystemHud />