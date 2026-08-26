package actiontype

import (
	"context"

	"master-control/model"
)

type Repository interface {
	Set(ctx context.Context, actionType model.ActionType) error
	GetAll(ctx context.Context) ([]model.ActionType, error)
	Get(ctx context.Context, ID string) (model.ActionType, error)
	Delete(ctx context.Context, ID string) error
}

type ActionType struct {
	repository Repository
}

func New(repository Repository) *ActionType {
	return &ActionType{
		repository: repository,
	}
}

func (a *ActionType) Set(ctx context.Context, actionType model.ActionType) error {
	return a.repository.Set(ctx, actionType)
}

func (a *ActionType) GetAll(ctx context.Context) ([]model.ActionType, error) {
	return a.repository.GetAll(ctx)
}

func (a *ActionType) Get(ctx context.Context, ID string) (model.ActionType, error) {
	return a.repository.Get(ctx, ID)
}

func (a *ActionType) Delete(ctx context.Context, ID string) error {
	return a.repository.Delete(ctx, ID)
}
