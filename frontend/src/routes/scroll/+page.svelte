<script lang="ts">
    import { onDestroy } from 'svelte';
    import { scrollBy } from '$lib/api/system';

    // ── Preferencias (por dispositivo) ────────────────────────────────────
    const PREFS_KEY = 'master-control:scroll';

    type Prefs = { speed: number; invert: boolean; horizontal: boolean; inertia: boolean };
    const defaults: Prefs = { speed: 2, invert: false, horizontal: true, inertia: true };

    function loadPrefs(): Prefs {
        try {
            const raw = localStorage.getItem(PREFS_KEY);
            if (raw) return { ...defaults, ...JSON.parse(raw) };
        } catch {
            // almacenamiento bloqueado o JSON roto: valores por defecto
        }
        return { ...defaults };
    }

    let prefs = $state<Prefs>(typeof window === 'undefined' ? { ...defaults } : loadPrefs());

    $effect(() => {
        const snapshot = JSON.stringify(prefs);
        try {
            localStorage.setItem(PREFS_KEY, snapshot);
        } catch {
            // sin persistencia: las preferencias duran esta sesión
        }
    });

    // ── Envío: un request en vuelo a la vez, acumulando lo que llega mientras tanto ──
    let pendingX = 0;
    let pendingY = 0;
    let inFlight = false;
    let error = $state('');
    let active = $state(false);

    function queue(dx: number, dy: number) {
        const sign = prefs.invert ? -1 : 1;
        // El contenido sigue al dedo (como una pantalla táctil): dedo hacia abajo → la
        // página baja → rueda positiva en CoreGraphics.
        pendingX += prefs.horizontal ? dx * prefs.speed * sign : 0;
        pendingY += dy * prefs.speed * sign;
        flush();
    }

    async function flush() {
        if (inFlight) return;
        const dx = Math.trunc(pendingX);
        const dy = Math.trunc(pendingY);
        if (dx === 0 && dy === 0) return;
        pendingX -= dx;
        pendingY -= dy;
        inFlight = true;
        try {
            const res = await scrollBy(dx, dy);
            if (res.ok) {
                error = '';
            } else {
                error = res.error;
                stopInertia();
                pendingX = pendingY = 0;
            }
        } catch {
            error = 'No se pudo contactar a la Mac';
            stopInertia();
            pendingX = pendingY = 0;
        } finally {
            inFlight = false;
        }
        flush();
    }

    // ── Gesto ─────────────────────────────────────────────────────────────
    let pointerId: number | null = null;
    let lastX = 0;
    let lastY = 0;
    let lastT = 0;
    let velX = 0; // px/ms, suavizada
    let velY = 0;

    function onDown(e: PointerEvent) {
        if (pointerId !== null) return; // un solo dedo manda
        stopInertia();
        pointerId = e.pointerId;
        (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
        lastX = e.clientX;
        lastY = e.clientY;
        lastT = e.timeStamp;
        velX = velY = 0;
        active = true;
    }

    function onMove(e: PointerEvent) {
        if (e.pointerId !== pointerId) return;
        const dx = e.clientX - lastX;
        const dy = e.clientY - lastY;
        const dt = Math.max(e.timeStamp - lastT, 1);
        lastX = e.clientX;
        lastY = e.clientY;
        lastT = e.timeStamp;
        velX = 0.7 * (dx / dt) + 0.3 * velX;
        velY = 0.7 * (dy / dt) + 0.3 * velY;
        queue(dx, dy);
    }

    function onUp(e: PointerEvent) {
        if (e.pointerId !== pointerId) return;
        pointerId = null;
        active = false;
        // Si el dedo se quedó quieto antes de soltar, no hay inercia.
        if (e.timeStamp - lastT > 80) velX = velY = 0;
        if (prefs.inertia) startInertia();
    }

    // ── Inercia: sigue desplazando y frena, como el trackpad ──────────────
    let raf = 0;

    function startInertia() {
        if (Math.hypot(velX, velY) < 0.05) return;
        let prev = performance.now();
        const step = (now: number) => {
            const dt = Math.min(now - prev, 50);
            prev = now;
            const decay = Math.pow(0.995, dt); // ~ -0.5 % por ms
            velX *= decay;
            velY *= decay;
            if (Math.hypot(velX, velY) < 0.02) {
                raf = 0;
                return;
            }
            queue(velX * dt, velY * dt);
            raf = requestAnimationFrame(step);
        };
        raf = requestAnimationFrame(step);
    }

    function stopInertia() {
        if (raf) cancelAnimationFrame(raf);
        raf = 0;
    }

    function page(direction: 1 | -1) {
        stopInertia();
        // Un "salto" de ~ media pantalla; no depende de la velocidad.
        pendingY += 600 * direction;
        flush();
    }

    onDestroy(stopInertia);
</script>

<div class="wrap">
    <div
        class="pad"
        class:active
        role="application"
        aria-label="Área de scroll: deslice el dedo para desplazar la Mac"
        onpointerdown={onDown}
        onpointermove={onMove}
        onpointerup={onUp}
        onpointercancel={onUp}
    >
        <span class="hint">Deslice para desplazar<br /><small>La Mac mueve la ventana que está bajo el puntero</small></span>
    </div>

    {#if error}
        <div class="error" role="alert">{error}</div>
    {/if}

    <div class="controls">
        <button type="button" class="step" onclick={() => page(1)} aria-label="Subir media pantalla">⇞</button>
        <button type="button" class="step" onclick={() => page(-1)} aria-label="Bajar media pantalla">⇟</button>

        <label class="speed">
            Velocidad <strong>{prefs.speed}×</strong>
            <input type="range" min="0.5" max="6" step="0.5" bind:value={prefs.speed} />
        </label>

        <label class="toggle"><input type="checkbox" bind:checked={prefs.invert} /> Invertir</label>
        <label class="toggle"><input type="checkbox" bind:checked={prefs.horizontal} /> Horizontal</label>
        <label class="toggle"><input type="checkbox" bind:checked={prefs.inertia} /> Inercia</label>
    </div>
</div>

<style>
    .wrap {
        display: flex;
        flex-direction: column;
        gap: 12px;
        height: 100%;
        padding: 16px;
        box-sizing: border-box;
    }

    .pad {
        flex: 1;
        min-height: 240px;
        display: grid;
        place-items: center;
        border-radius: 20px;
        background: #1a1a1a;
        border: 1px solid rgba(255, 255, 255, 0.08);
        /* El gesto es nuestro: sin scroll, zoom ni menú del navegador. */
        touch-action: none;
        user-select: none;
        -webkit-user-select: none;
        -webkit-touch-callout: none;
        transition: background 120ms ease, border-color 120ms ease;
    }

    .pad.active {
        background: #222;
        border-color: rgba(255, 255, 255, 0.2);
    }

    .hint {
        text-align: center;
        color: rgba(255, 255, 255, 0.45);
        font-size: 1.1rem;
        line-height: 1.6;
        pointer-events: none;
    }

    .hint small {
        font-size: 0.8rem;
        color: rgba(255, 255, 255, 0.3);
    }

    .error {
        padding: 10px 14px;
        border-radius: 12px;
        background: rgba(255, 80, 80, 0.12);
        border: 1px solid rgba(255, 80, 80, 0.35);
        color: #ff9d9d;
        font-size: 0.9rem;
    }

    .controls {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 10px 16px;
        color: rgba(255, 255, 255, 0.8);
        font-size: 0.95rem;
    }

    .step {
        width: 56px;
        height: 48px;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.12);
        background: #1f1f1f;
        color: white;
        font-size: 1.4rem;
        cursor: pointer;
        touch-action: manipulation;
    }

    .step:active {
        background: #2c2c2c;
    }

    .speed {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .speed input {
        width: 140px;
    }

    .toggle {
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .toggle input {
        width: 20px;
        height: 20px;
    }
</style>
