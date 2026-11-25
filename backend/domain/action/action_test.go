package action

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"master-control/model"
)

type fakeRepo struct{ actions map[string]model.Action }

func (f fakeRepo) Set(context.Context, model.Action) error { return nil }
func (f fakeRepo) GetAll(context.Context) ([]model.Action, error) {
	return nil, nil
}
func (f fakeRepo) Get(_ context.Context, id string) (model.Action, error) {
	a, ok := f.actions[id]
	if !ok {
		return model.Action{}, errors.New("not found")
	}
	return a, nil
}
func (f fakeRepo) Delete(context.Context, string) error { return nil }

type fakeTypes struct{}

func (fakeTypes) Set(context.Context, model.ActionType) error { return nil }
func (fakeTypes) GetAll(context.Context) ([]model.ActionType, error) {
	return nil, nil
}
func (fakeTypes) Get(_ context.Context, id string) (model.ActionType, error) {
	if id == "command" {
		return model.ActionType{ID: "command", Key: "command", CommandTemplate: "bash -c"}, nil
	}
	return model.ActionType{}, errors.New("unknown type")
}
func (fakeTypes) Delete(context.Context, string) error { return nil }

type fakeSystem struct{ called string }

func (f *fakeSystem) Execute(_ context.Context, token string) error {
	f.called = token
	return nil
}

func newUseCase(maxTimeout time.Duration, actions ...model.Action) (*Action, *fakeSystem) {
	m := map[string]model.Action{}
	for _, a := range actions {
		m[a.ID] = a
	}
	sys := &fakeSystem{}
	return New(fakeRepo{actions: m}, fakeTypes{}, sys, maxTimeout), sys
}

func TestExecuteCapturaSalida(t *testing.T) {
	uc, _ := newUseCase(0, model.Action{ID: "a", Type: "command", Command: "echo hola; echo fallo >&2"})
	res, err := uc.Execute(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK || res.ExitCode != 0 || strings.TrimSpace(res.Stdout) != "hola" || strings.TrimSpace(res.Stderr) != "fallo" {
		t.Errorf("resultado inesperado: %+v", res)
	}
}

func TestExecuteExitCodeNoEsError(t *testing.T) {
	uc, _ := newUseCase(0, model.Action{ID: "a", Type: "command", Command: "exit 3"})
	res, err := uc.Execute(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	if res.OK || res.ExitCode != 3 {
		t.Errorf("quiero OK=false exit=3, got %+v", res)
	}
}

func TestExecuteTimeout(t *testing.T) {
	uc, _ := newUseCase(time.Second, model.Action{ID: "a", Type: "command", Command: "sleep 5", TimeoutSeconds: 30})
	start := time.Now()
	res, err := uc.Execute(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	if !res.TimedOut || res.OK || res.ExitCode != -1 {
		t.Errorf("quiero timeout, got %+v", res)
	}
	if time.Since(start) > 4*time.Second {
		t.Errorf("el tope MC_MAX_TIMEOUT no se aplicó: %v", time.Since(start))
	}
}

func TestExecuteTrunca(t *testing.T) {
	uc, _ := newUseCase(0, model.Action{ID: "a", Type: "command", Command: "head -c 40000 /dev/zero | tr '\\0' x"})
	res, err := uc.Execute(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Truncated || len(res.Stdout) != maxOutputBytes || !res.OK {
		t.Errorf("quiero salida truncada a %d, got len=%d truncated=%v ok=%v", maxOutputBytes, len(res.Stdout), res.Truncated, res.OK)
	}
}

func TestExecuteSystem(t *testing.T) {
	uc, sys := newUseCase(0, model.Action{ID: "v", Type: "system", Command: "volume_up"})
	res, err := uc.Execute(context.Background(), "v")
	if err != nil || !res.OK || sys.called != "volume_up" {
		t.Errorf("res=%+v err=%v called=%q", res, err, sys.called)
	}
}

func TestExecuteErrores(t *testing.T) {
	uc, _ := newUseCase(0, model.Action{ID: "x", Type: "desconocido", Command: "true"})
	if _, err := uc.Execute(context.Background(), "nope"); err == nil {
		t.Error("acción inexistente debe ser error")
	}
	if _, err := uc.Execute(context.Background(), "x"); err == nil {
		t.Error("tipo inexistente debe ser error")
	}
}

func TestTimeoutFor(t *testing.T) {
	uc, _ := newUseCase(20 * time.Second)
	if got := uc.timeoutFor(model.Action{}); got != defaultTimeout {
		t.Errorf("default = %v", got)
	}
	if got := uc.timeoutFor(model.Action{TimeoutSeconds: 15}); got != 15*time.Second {
		t.Errorf("15s = %v", got)
	}
	if got := uc.timeoutFor(model.Action{TimeoutSeconds: 90}); got != 20*time.Second {
		t.Errorf("recortado = %v", got)
	}
}
