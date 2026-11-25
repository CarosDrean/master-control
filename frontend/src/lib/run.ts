import { executeAction, type Action, type ExecutionResult } from '$lib/api/actions';
import { executeSystemActionWithFeedback } from '$lib/api/system';

// Flujo único para ejecutar una acción desde la UI:
// confirmación (si la acción la pide) → ejecución → feedback (aviso breve o salida completa).
// El componente ActionFeedback (montado en el layout) registra la vista.

export interface RunFeedback {
    action: Action;
    result: ExecutionResult;
    /** true = abrir el panel con la salida; false = solo un aviso breve. */
    expanded: boolean;
}

interface FeedbackView {
    confirm(action: Action): Promise<boolean>;
    running(action: Action, on: boolean): void;
    show(feedback: RunFeedback): void;
}

let view: FeedbackView | null = null;

export function registerFeedback(v: FeedbackView): () => void {
    view = v;
    return () => {
        if (view === v) view = null;
    };
}

/** Ejecuta la acción con su confirmación y feedback. Devuelve null si se canceló. */
export async function runAction(action: Action): Promise<ExecutionResult | null> {
    if (action.confirm && view && !(await view.confirm(action))) {
        return null;
    }

    view?.running(action, true);
    let result: ExecutionResult;
    try {
        result = action.type === "system"
            ? await executeSystemActionWithFeedback(action)
            : await executeAction(action.id);
    } finally {
        view?.running(action, false);
    }

    // El HUD de sistema ya da feedback propio; el resto avisa siempre, y un fallo abre la salida.
    const silent = action.type === "system" && result.ok;
    if (!silent) {
        view?.show({ action, result, expanded: !!action.show_output || !result.ok });
    }

    return result;
}
