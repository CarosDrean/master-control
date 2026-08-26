package action

import (
	"context"

	"master-control/model"
)

type Repository interface {
	Set(ctx context.Context, action model.Action) error
	GetAll(ctx context.Context) ([]model.Action, error)
	Get(ctx context.Context, id string) (model.Action, error)
	Delete(ctx context.Context, id string) error
}

type UseCaseActionType interface {
	Set(ctx context.Context, actionType model.ActionType) error
	GetAll(ctx context.Context) ([]model.ActionType, error)
	Get(ctx context.Context, ID string) (model.ActionType, error)
	Delete(ctx context.Context, ID string) error
}

type UseCaseSystem interface {
	Execute(ctx context.Context, token string) error
}
