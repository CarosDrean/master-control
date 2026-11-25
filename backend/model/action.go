package model

type Action struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Icon     string `json:"icon"`
	Type     string `json:"type"`
	Command  string `json:"command"`
	Page     string `json:"page"`
	Position int    `json:"position"`
	// Confirm pide confirmación en la UI antes de ejecutar (acciones delicadas).
	Confirm bool `json:"confirm,omitempty"`
	// ShowOutput muestra la salida del comando en la UI al terminar.
	ShowOutput bool `json:"show_output,omitempty"`
	// TimeoutSeconds acota la ejecución; 0 = timeout por defecto. Se recorta a MC_MAX_TIMEOUT.
	TimeoutSeconds int `json:"timeout_seconds,omitempty"`
}

type Actions []Action

// ExecutionResult es el resultado de ejecutar una acción.
type ExecutionResult struct {
	OK         bool   `json:"ok"`
	ExitCode   int    `json:"exit_code"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"duration_ms"`
	TimedOut   bool   `json:"timed_out"`
	Truncated  bool   `json:"truncated"`
}
