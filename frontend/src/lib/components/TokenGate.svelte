<script lang="ts">
    import { onMount, type Snippet } from 'svelte';
    import { adoptTokenFromUrl, getToken, onAuthRequired, setToken, verifyToken } from '$lib/api/api';

    let { children }: { children?: Snippet } = $props();

    // checking → ok | locked. Sin token válido no se monta la app (evita una cascada de 401).
    let phase = $state<'checking' | 'ok' | 'locked'>('checking');
    let input = $state('');
    let error = $state('');

    async function check() {
        phase = 'checking';
        phase = getToken() && (await verifyToken()) ? 'ok' : 'locked';
    }

    async function submit(e: Event) {
        e.preventDefault();
        setToken(input);
        await check();
        error = phase === 'locked' ? 'Token incorrecto o servidor no disponible.' : '';
        if (phase === 'ok') input = '';
    }

    onMount(() => {
        adoptTokenFromUrl();
        check();
        return onAuthRequired(() => {
            if (phase === 'ok') {
                error = 'El servidor rechazó el token.';
                phase = 'locked';
            }
        });
    });
</script>

{#if phase === 'ok'}
    {@render children?.()}
{:else if phase === 'checking'}
    <div class="center muted">Conectando…</div>
{:else}
    <form class="center" onsubmit={submit}>
        <h2>Token de acceso</h2>
        <p class="muted">
            Ábrelo desde la URL que imprime <code>./dev.sh</code>, o pega el valor de <code>MC_TOKEN</code>
            (en <code>backend/.env</code>).
        </p>
        <input
            type="password"
            bind:value={input}
            placeholder="MC_TOKEN"
            autocomplete="current-password"
            autocapitalize="off"
            spellcheck="false"
        />
        {#if error}<p class="error">{error}</p>{/if}
        <button type="submit" disabled={!input.trim()}>Entrar</button>
    </form>
{/if}

<style>
    .center {
        max-width: 420px;
        margin: 12vh auto 0;
        padding: 0 16px;
        display: flex;
        flex-direction: column;
        gap: 14px;
        color: white;
        text-align: center;
    }

    h2 {
        margin: 0;
    }

    .muted {
        color: rgba(255, 255, 255, 0.55);
        font-size: 0.92rem;
        line-height: 1.5;
        margin: 0;
    }

    code {
        background: rgba(255, 255, 255, 0.08);
        padding: 1px 6px;
        border-radius: 6px;
    }

    input {
        background: #1f1f1f;
        border: 1px solid rgba(255, 255, 255, 0.15);
        border-radius: 12px;
        padding: 14px;
        color: white;
        font-size: 1rem;
        outline: none;
    }

    button {
        min-height: 52px;
        border-radius: 12px;
        border: none;
        background: #2a5bd7;
        color: white;
        font-size: 1rem;
        cursor: pointer;
    }

    button:disabled {
        opacity: 0.5;
        cursor: default;
    }

    .error {
        color: #ffb3b3;
        margin: 0;
    }
</style>
