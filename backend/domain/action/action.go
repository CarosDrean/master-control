package action

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"master-control/model"
)

// Timeout por defecto de una acción y tope de salida capturada por stream.
const (
	defaultTimeout = 10 * time.Second
	maxOutputBytes = 16 * 1024
)

type Action struct {
	repository Repository
	actionType UseCaseActionType
	system     UseCaseSystem
	maxTimeout time.Duration
}

// New crea el caso de uso. maxTimeout acota el timeout_seconds de cada acción
// (0 = sin tope adicional al timeout por defecto).
func New(repository Repository, actionType UseCaseActionType, system UseCaseSystem, maxTimeout time.Duration) *Action {
	return &Action{
		repository: repository,
		actionType: actionType,
		system:     system,
		maxTimeout: maxTimeout,
	}
}

// timeoutFor resuelve el timeout efectivo de una acción.
func (a *Action) timeoutFor(action model.Action) time.Duration {
	timeout := defaultTimeout
	if action.TimeoutSeconds > 0 {
		timeout = time.Duration(action.TimeoutSeconds) * time.Second
	}
	if a.maxTimeout > 0 && timeout > a.maxTimeout {
		timeout = a.maxTimeout
	}
	return timeout
}

// Execute ejecuta la acción y devuelve su salida. Un exit code distinto de cero
// o un timeout se reportan en el resultado (OK=false), no como error: el error
// queda para lo que impide ejecutar (acción o tipo inexistente, binario ausente).
func (a *Action) Execute(ctx context.Context, id string) (model.ExecutionResult, error) {
	action, err := a.repository.Get(ctx, id)
	if err != nil {
		return model.ExecutionResult{}, err
	}

	start := time.Now()

	if action.Type == "system" {
		if err := a.system.Execute(ctx, action.Command); err != nil {
			return model.ExecutionResult{}, err
		}
		return model.ExecutionResult{OK: true, DurationMs: time.Since(start).Milliseconds()}, nil
	}

	actionType, err := a.actionType.Get(ctx, action.Type)
	if err != nil {
		return model.ExecutionResult{}, fmt.Errorf("actionType.Get(): %w", err)
	}

	commands := strings.Fields(actionType.CommandTemplate)
	if len(commands) == 0 {
		return model.ExecutionResult{}, fmt.Errorf("invalid command template for action type: %s", actionType.Key)
	}

	executable := commands[0]
	args := append(commands[1:], action.Command)

	execCtx, cancel := context.WithTimeout(ctx, a.timeoutFor(action))
	defer cancel()

	stdout := &limitedBuffer{limit: maxOutputBytes}
	stderr := &limitedBuffer{limit: maxOutputBytes}

	cmd := exec.CommandContext(execCtx, executable, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	// Un hijo que herede stdout no debe colgar la respuesta tras el timeout.
	cmd.WaitDelay = 2 * time.Second

	runErr := cmd.Run()

	result := model.ExecutionResult{
		OK:         runErr == nil,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: time.Since(start).Milliseconds(),
		TimedOut:   errors.Is(execCtx.Err(), context.DeadlineExceeded),
		Truncated:  stdout.truncated || stderr.truncated,
	}

	if runErr != nil {
		var exitErr *exec.ExitError
		switch {
		case result.TimedOut:
			result.ExitCode = -1
		case errors.As(runErr, &exitErr):
			result.ExitCode = exitErr.ExitCode()
		default:
			return model.ExecutionResult{}, runErr
		}
	}

	return result, nil
}

// limitedBuffer guarda hasta limit bytes y descarta el resto, marcando truncated.
type limitedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		if len(p) > room {
			b.buf.Write(p[:room])
			b.truncated = true
		} else {
			b.buf.Write(p)
		}
	} else if len(p) > 0 {
		b.truncated = true
	}
	// Se reporta todo como escrito para no cortar el proceso con un error de pipe.
	return len(p), nil
}

func (b *limitedBuffer) String() string {
	return b.buf.String()
}

func (a *Action) Set(ctx context.Context, action model.Action) error {
	return a.repository.Set(ctx, action)
}

func (a *Action) GetAll(ctx context.Context) ([]model.Action, error) {
	return a.repository.GetAll(ctx)
}

func (a *Action) Get(ctx context.Context, id string) (model.Action, error) {
	return a.repository.Get(ctx, id)
}

func (a *Action) Delete(ctx context.Context, id string) error {
	return a.repository.Delete(ctx, id)
}
