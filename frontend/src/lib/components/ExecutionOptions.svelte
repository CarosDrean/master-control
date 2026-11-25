<script lang="ts">
    // Opciones de ejecución de una acción (compartido por crear y editar).
    export let type = "";
    export let confirm = false;
    export let showOutput = false;
    export let timeoutStr = "";

    // Las acciones de sistema no tienen salida ni timeout propio.
    $: hasProcess = type !== "system";
</script>

<fieldset class="exec">
    <legend>Ejecución</legend>

    <label class="check">
        <input type="checkbox" bind:checked={confirm} />
        <span>Pedir confirmación antes de ejecutar</span>
    </label>

    {#if hasProcess}
        <label class="check">
            <input type="checkbox" bind:checked={showOutput} />
            <span>Mostrar la salida al terminar <small>(un fallo la muestra siempre)</small></span>
        </label>

        <label class="timeout" for="timeout-seconds">
            <span>Timeout (segundos, vacío = 10)</span>
            <input id="timeout-seconds" bind:value={timeoutStr} inputmode="numeric" placeholder="10" />
        </label>
    {/if}
</fieldset>

<style>
    .exec {
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        padding: 12px 14px 4px;
        margin: 0 0 16px;
    }

    legend {
        font-size: 0.85rem;
        opacity: 0.8;
        padding: 0 6px;
    }

    .check {
        display: flex;
        align-items: center;
        gap: 12px;
        min-height: 44px;
        cursor: pointer;
        font-size: 0.95rem;
    }

    .check input {
        width: 22px;
        height: 22px;
        flex: 0 0 auto;
        accent-color: #2a5bd7;
    }

    small {
        opacity: 0.55;
    }

    .timeout {
        display: flex;
        flex-direction: column;
        gap: 6px;
        margin: 6px 0 12px;
        font-size: 0.85rem;
    }

    .timeout span {
        opacity: 0.8;
    }

    .timeout input {
        background: #1f1f1f;
        border: 1px solid rgba(255, 255, 255, 0.12);
        border-radius: 10px;
        padding: 12px 14px;
        color: white;
        font-size: 0.95rem;
        max-width: 140px;
        outline: none;
    }
</style>
