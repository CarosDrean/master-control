// Mantiene la pantalla encendida mientras Master Control está visible (el iPad como control remoto
// no debe bloquearse ni atenuarse a mitad de uso).
//
// - Contexto seguro (localhost / HTTPS): Screen Wake Lock API. El sistema la suelta al ocultar la
//   pestaña, así que se vuelve a pedir al volver a ser visible.
// - Sin contexto seguro (el iPad por http en la red local): Safari no expone la API; un video mudo
//   y diminuto en bucle (`static/keepawake.mp4`) evita el bloqueo automático de iOS.
//
// Ambos necesitan un gesto del usuario la primera vez: arranca en el primer toque.

let sentinel: WakeLockSentinel | null = null;
let video: HTMLVideoElement | null = null;

function hasWakeLock(): boolean {
    return window.isSecureContext && "wakeLock" in navigator;
}

async function acquire(): Promise<void> {
    if (document.visibilityState !== "visible") return;
    if (hasWakeLock()) {
        if (sentinel && !sentinel.released) return;
        try {
            sentinel = await navigator.wakeLock.request("screen");
        } catch {
            // denegado (batería baja, política del navegador): se reintenta en el próximo toque
        }
        return;
    }
    if (!video) {
        video = document.createElement("video");
        video.src = "/keepawake.mp4";
        video.muted = true;
        video.loop = true;
        video.playsInline = true;
        video.setAttribute("playsinline", "");
        video.setAttribute("aria-hidden", "true");
        // En el DOM pero invisible: iOS no cuenta un video fuera del documento.
        Object.assign(video.style, {
            position: "fixed", width: "1px", height: "1px", opacity: "0", pointerEvents: "none", left: "0", top: "0",
        });
        document.body.appendChild(video);
    }
    if (video.paused) {
        video.play().catch(() => {
            // sin gesto todavía: se reintenta en el próximo toque
        });
    }
}

/** Activa el bloqueo de reposo; devuelve la función que lo desactiva. */
export function keepAwake(): () => void {
    const onGesture = () => void acquire();
    const onVisibility = () => void acquire();

    document.addEventListener("pointerdown", onGesture, { passive: true });
    document.addEventListener("visibilitychange", onVisibility);
    void acquire();

    return () => {
        document.removeEventListener("pointerdown", onGesture);
        document.removeEventListener("visibilitychange", onVisibility);
        void sentinel?.release().catch(() => {});
        sentinel = null;
        video?.pause();
        video?.remove();
        video = null;
    };
}
