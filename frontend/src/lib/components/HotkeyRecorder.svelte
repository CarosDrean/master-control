<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    import { captureHotkey } from '$lib/api/hotkeys';

    export let command = "";
    export let onChange: (value: string) => void = () => {};
    export let onCapture: (info: { modifiers: string[]; key: string }) => void = () => {};
    export let id = "hotkey-command";

    let recording = false;
    let capturingNative = false;
    let heldMods: string[] = [];
    let lastCombo: string[] = [];
    let modsEverPressed = false;
    let warning = "";
    let manual = "";
    let manualError = "";

    const MODIFIER_KEYS = new Set(["meta", "control", "alt", "shift"]);

    const PLACEHOLDER = 'tell application "System Events"\n\tkeystroke "p" using {command down}\nend tell';

    const MOD_LABELS: Record<string, string> = {
        meta: "⌘",
        control: "⌃",
        alt: "⌥",
        shift: "⇧",
    };

    const KEY_CODES: Record<string, number> = {
        Enter: 36,
        Tab: 48,
        Space: 49,
        Backspace: 51,
        Escape: 53,
        ArrowLeft: 123,
        ArrowRight: 124,
        ArrowDown: 125,
        ArrowUp: 126,
        Home: 115,
        End: 119,
        PageUp: 116,
        PageDown: 121,
        F1: 122, F2: 120, F3: 99, F4: 118, F5: 96, F6: 97, F7: 98, F8: 100,
        F9: 101, F10: 109, F11: 103, F12: 111, F13: 105, F14: 107, F15: 113,
        F16: 106, F17: 64, F18: 79, F19: 80, F20: 90,
    };

    const CODE_CHARS: Record<string, string> = {
        Space: " ", Minus: "-", Equal: "=", BracketLeft: "[", BracketRight: "]",
        Backslash: "\\", Semicolon: ";", Quote: "'", Comma: ",", Period: ".",
        Slash: "/", Backquote: "`", IntlBackslash: "\\",
    };

    const KEYCODE_LABELS: Record<number, string> = {
        0: "A", 1: "S", 2: "D", 3: "F", 4: "H", 5: "G", 6: "Z", 7: "X",
        8: "C", 9: "V", 11: "B", 12: "Q", 13: "W", 14: "E", 15: "R",
        16: "Y", 17: "T", 18: "1", 19: "2", 20: "3", 21: "4", 22: "6",
        23: "5", 24: "=", 25: "9", 26: "7", 27: "-", 28: "8", 29: "0",
        30: "]", 31: "O", 32: "U", 33: "[", 34: "I", 35: "P", 37: "L",
        38: "J", 39: "'", 40: "K", 41: ";", 42: "\\", 43: ",", 44: "/",
        45: "N", 46: "M", 47: ".", 48: "Tab", 49: "Space", 50: "`",
        51: "Del", 53: "Esc", 65: ".", 67: "*", 69: "+", 71: "Clear",
        75: "/", 76: "Enter", 78: "-", 82: "0", 83: "1", 84: "2", 85: "3",
        86: "4", 87: "5", 88: "6", 89: "7", 91: "8", 92: "9",
        96: "F5", 97: "F6", 98: "F7", 99: "F3", 100: "F8", 101: "F9",
        103: "F11", 105: "F13", 106: "F16", 107: "F14", 109: "F10",
        111: "F12", 113: "F15", 114: "Insert", 115: "Home", 116: "PgUp",
        117: "FwdDel", 119: "End", 120: "F2", 121: "PgDn", 122: "F1",
        123: "Left", 124: "Right", 125: "Down", 126: "Up",
    };

    const MOD_ALIASES: Record<string, string> = {
        cmd: "meta", command: "meta", meta: "meta",
        ctrl: "control", control: "control",
        opt: "alt", option: "alt", alt: "alt",
        shift: "shift",
    };

    const SPECIAL_KEYS: Record<string, string> = {
        space: "Space", enter: "Enter", return: "Enter",
        tab: "Tab", esc: "Escape", escape: "Escape",
        backspace: "Backspace", delete: "Backspace",
        home: "Home", end: "End", pageup: "PageUp", pagedown: "PageDown",
        up: "ArrowUp", down: "ArrowDown", left: "ArrowLeft", right: "ArrowRight",
    };

    function keyOfCode(code: string): string | null {
        if (/^Key[A-Z]$/.test(code)) return code.slice(3).toLowerCase();
        if (/^Digit[0-9]$/.test(code)) return code.slice(5);
        if (KEY_CODES[code] !== undefined) return code;
        return CODE_CHARS[code] ?? null;
    }

    function slugify(s: string): string {
        return s
            .toLowerCase()
            .normalize("NFD")
            .replace(/[\u0300-\u036f]/g, "")
            .replace(/[^a-z0-9]+/g, "");
    }

    function slugKey(key: string): string {
        const map: Record<string, string> = {
            " ": "space", Enter: "enter", Tab: "tab", Backspace: "backspace",
            Escape: "esc", ArrowLeft: "left", ArrowRight: "right",
            ArrowUp: "up", ArrowDown: "down", Home: "home", End: "end",
            PageUp: "pageup", PageDown: "pagedown",
        };
        if (map[key]) return map[key];
        const slug = slugify(key);
        return slug || "key";
    }

    function keycodeSlug(keycode: number): string {
        const slug = slugify(keycodeLabel(keycode));
        return slug || String(keycode);
    }

    function modsFromEvent(e: KeyboardEvent): string[] {
        const m: string[] = [];
        if (e.metaKey) m.push("meta");
        if (e.ctrlKey) m.push("control");
        if (e.altKey) m.push("alt");
        if (e.shiftKey) m.push("shift");
        return m;
    }

    function clearWarning() {
        warning = "";
    }

    function startRecording() {
        recording = true;
        heldMods = [];
        lastCombo = [];
        modsEverPressed = false;
        clearWarning();
    }

    function stopRecording() {
        recording = false;
        heldMods = [];
        modsEverPressed = false;
    }

    function buildCommand(mods: string[], key: string): string {
        const using: string[] = [];
        if (mods.includes("shift")) using.push("shift down");
        if (mods.includes("control")) using.push("control down");
        if (mods.includes("alt")) using.push("option down");
        if (mods.includes("meta")) using.push("command down");
        const clause = using.length > 0 ? ` using {${using.join(", ")}}` : "";

        const keyCode = KEY_CODES[key];
        const line = keyCode !== undefined
            ? `key code ${keyCode}${clause}`
            : `keystroke "${key}"${clause}`;

        return `tell application "System Events"\n\t${line}\nend tell`;
    }

    function setCombo(mods: string[], key: string) {
        lastCombo = [...mods.map(m => MOD_LABELS[m]), key];
        const cmd = buildCommand(mods, key);
        command = cmd;
        onChange(cmd);
        onCapture({ modifiers: mods, key: slugKey(key) });
        clearWarning();
    }

    function onKeyDown(e: KeyboardEvent) {
        if (!recording) return;
        e.preventDefault();
        e.stopPropagation();

        if (MODIFIER_KEYS.has(e.key)) {
            modsEverPressed = true;
            heldMods = modsFromEvent(e);
            return;
        }

        if (e.key === "Escape") {
            stopRecording();
            return;
        }

        const key = keyOfCode(e.code);
        if (key === null) return;

        setCombo(modsFromEvent(e), key);
        stopRecording();
    }

    function onKeyUp(e: KeyboardEvent) {
        if (!recording) return;

        if (MODIFIER_KEYS.has(e.key)) {
            heldMods = modsFromEvent(e);
            if (heldMods.length === 0 && modsEverPressed) {
                warning = "Parece que el sistema interceptó el atajo (Spotlight, Mission Control, ⌘Tab…). El navegador no puede capturarlo: usa el campo manual de abajo.";
                stopRecording();
            }
        }
    }

    function onBlur() {
        if (!recording) return;
        warning = "La ventana perdió el foco, así que el sistema interceptó el atajo. Usa el campo manual de abajo.";
        stopRecording();
    }

    function parseManual(input: string): boolean {
        const parts = input.toLowerCase().split(/[\s+,]+/).filter(Boolean);
        const mods: string[] = [];
        let key = "";

        for (const p of parts) {
            if (MOD_ALIASES[p]) {
                mods.push(MOD_ALIASES[p]);
            } else if (SPECIAL_KEYS[p]) {
                key = SPECIAL_KEYS[p];
            } else if (/^f([1-9]|1[0-9]|2[0])$/i.test(p)) {
                key = p.charAt(0).toUpperCase() + p.slice(1);
            } else if (p.length === 1) {
                key = p;
            } else {
                return false;
            }
        }

        if (!key) return false;

        setCombo(mods, key);
        return true;
    }

    function applyManual() {
        if (!parseManual(manual)) {
            manualError = "Formato inválido. Ejemplos: cmd+shift+p · option+space · ctrl+f12";
        } else {
            manualError = "";
        }
    }

    function keycodeLabel(keycode: number): string {
        return KEYCODE_LABELS[keycode] ?? `#${keycode}`;
    }

    async function captureNative() {
        if (capturingNative) return;
        capturingNative = true;
        clearWarning();
        manualError = "";

        const res = await captureHotkey();

        if (res.success && res.hotkey) {
            const h = res.hotkey;
            lastCombo = [...h.modifiers.map(m => MOD_LABELS[m] ?? m), keycodeLabel(h.keycode)];
            command = h.command;
            onChange(h.command);
            onCapture({ modifiers: h.modifiers, key: keycodeSlug(h.keycode) });
        } else {
            const status = res.status ?? 0;
            if (status === 403) {
                warning = "El backend no tiene permiso de Accesibilidad. Actívalo en Ajustes del Sistema > Privacidad y seguridad > Accesibilidad y vuelve a intentarlo.";
            } else if (status === 504) {
                warning = "Se agotó el tiempo de espera. Pulsa la combinación antes de 15 segundos.";
            } else {
                warning = res.message ?? "No se pudo capturar el hotkey.";
            }
        }

        capturingNative = false;
    }

    onMount(() => {
        window.addEventListener("keydown", onKeyDown, true);
        window.addEventListener("keyup", onKeyUp, true);
        window.addEventListener("blur", onBlur);
    });

    onDestroy(() => {
        window.removeEventListener("keydown", onKeyDown, true);
        window.removeEventListener("keyup", onKeyUp, true);
        window.removeEventListener("blur", onBlur);
    });

    $: comboDisplay = recording
        ? [...heldMods.map(m => MOD_LABELS[m])]
        : lastCombo;
</script>

<style>
    .recorder {
        display: flex;
        align-items: center;
        gap: 14px;
        margin-bottom: 12px;
    }

    .btn-record {
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 12px;
        padding: 10px 18px;
        font-size: 0.92rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
        white-space: nowrap;
    }

    .btn-record:hover {
        background: #292929;
    }

    .btn-record.recording {
        background: rgba(255, 80, 80, 0.15);
        border-color: rgba(255, 80, 80, 0.5);
        color: #ff8080;
        animation: pulse 1.2s ease infinite;
    }

    .btn-native {
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 12px;
        padding: 10px 16px;
        font-size: 0.88rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
        white-space: nowrap;
        margin-bottom: 12px;
    }

    .btn-native:hover {
        background: #292929;
    }

    .btn-native:disabled {
        opacity: 0.6;
        cursor: wait;
        color: #ffb3b3;
        border-color: rgba(255, 80, 80, 0.5);
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.55; }
    }

    .combo {
        display: flex;
        align-items: center;
        gap: 6px;
        min-height: 38px;
    }

    .key {
        min-width: 30px;
        height: 30px;
        padding: 0 8px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.18);
        border-bottom-width: 2px;
        border-radius: 8px;
        font-size: 0.9rem;
        color: white;
        box-shadow: 0 3px 6px rgba(0,0,0,0.35);
    }

    .key.placeholder {
        opacity: 0.4;
        font-size: 0.8rem;
        border-style: dashed;
        box-shadow: none;
    }

    .hint {
        font-size: 0.8rem;
        opacity: 0.5;
        margin-bottom: 10px;
    }

    .warning {
        margin: 0 0 10px;
        padding: 10px 12px;
        font-size: 0.82rem;
        line-height: 1.4;
        color: #ffb3b3;
        background: rgba(255, 80, 80, 0.12);
        border: 1px solid rgba(255, 80, 80, 0.4);
        border-radius: 10px;
    }

    .manual {
        display: flex;
        gap: 8px;
        margin-bottom: 12px;
    }

    .manual input {
        flex: 1;
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 10px 14px;
        color: white;
        font-size: 0.9rem;
        outline: none;
        min-width: 0;
    }

    .manual input:focus {
        border-color: rgba(255,255,255,0.25);
    }

    .btn-manual {
        background: #1e1e1e;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 10px 16px;
        font-size: 0.9rem;
        color: white;
        cursor: pointer;
        transition: 0.15s ease;
        white-space: nowrap;
    }

    .btn-manual:hover {
        background: #292929;
    }

    .manual-error {
        font-size: 0.8rem;
        color: #ffb3b3;
        margin: 0 0 12px;
    }

    textarea {
        background: #1f1f1f;
        border: 1px solid rgba(255,255,255,0.12);
        border-radius: 10px;
        padding: 12px 14px;
        color: white;
        font-size: 0.95rem;
        outline: none;
        min-height: 120px;
        resize: vertical;
        width: 100%;
        box-sizing: border-box;
        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    }

    textarea:focus {
        border-color: rgba(255,255,255,0.25);
    }
</style>

<div class="recorder">
    <button
        type="button"
        class="btn-record {recording ? 'recording' : ''}"
        on:click={recording ? stopRecording : startRecording}
    >
        {recording ? '● Grabando…' : '⌨ Grabar hotkey'}
    </button>

    <div class="combo">
        {#if comboDisplay.length > 0}
            {#each comboDisplay as part, i}
                <span class="key">{part}</span>
            {/each}
        {:else}
            <span class="key placeholder">?</span>
        {/if}
    </div>
</div>

{#if recording}
    <div class="hint">Pulsa la combinación de teclas · Esc para cancelar</div>
{:else if capturingNative}
    <div class="hint">Escuchando en todo el sistema… pulsa la combinación en cualquier app (funciona con atajos del sistema como ⌘Espacio o ⌘Tab)</div>
{:else}
    <div class="hint">Los atajos reservados del sistema (⌘Espacio, ⌘Tab…) no se pueden grabar en el navegador. Usa el botón de captura del sistema o el campo manual.</div>
{/if}

<button
    type="button"
    class="btn-native"
    disabled={capturingNative}
    on:click={captureNative}
>
    {capturingNative ? '● Escuchando en el sistema…' : '🎙 Capturar en todo el sistema'}
</button>

{#if warning}
    <div class="warning">{warning}</div>
{/if}

<div class="manual">
    <input
        bind:value={manual}
        on:keydown={(e) => e.key === "Enter" && applyManual()}
        placeholder="Manual: cmd+shift+p, option+space, ctrl+f12"
    />
    <button type="button" class="btn-manual" on:click={applyManual}>Aplicar</button>
</div>

{#if manualError}
    <div class="manual-error">{manualError}</div>
{/if}

<textarea
    id={id}
    bind:value={command}
    on:input={() => onChange(command)}
    placeholder={PLACEHOLDER}
></textarea>