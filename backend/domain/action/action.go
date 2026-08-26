package action

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"master-control/model"
)

const commandTimeout = 10 * time.Second

type Action struct {
	repository Repository
	actionType UseCaseActionType
	system     UseCaseSystem
}

func New(repository Repository, actionType UseCaseActionType, system UseCaseSystem) *Action {
	return &Action{
		repository: repository,
		actionType: actionType,
		system:     system,
	}
}

func (a *Action) Execute(ctx context.Context, id string) error {
	action, err := a.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if action.Type == "system" {
		return a.system.Execute(ctx, action.Command)
	}

	actionType, err := a.actionType.Get(ctx, action.Type)
	if err != nil {
		return fmt.Errorf("actionType.Get(): %w", err)
	}

	commands := strings.Split(actionType.CommandTemplate, " ")
	if len(commands) == 0 {
		return fmt.Errorf("invalid command template for action type: %s", actionType.Key)
	}

	executable := commands[0]
	args := commands[1:]
	args = append(args, action.Command)

	execCtx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(execCtx, executable, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
