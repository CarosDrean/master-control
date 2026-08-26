package model

import "context"

type ActionUseCase interface {
	Execute(ctx context.Context, id string) error
	Set(ctx context.Context, action Action) error
	GetAll(ctx context.Context) ([]Action, error)
	Get(ctx context.Context, id string) (Action, error)
	Delete(ctx context.Context, id string) error
}

type ActionTypeUseCase interface {
	Set(ctx context.Context, actionType ActionType) error
	GetAll(ctx context.Context) ([]ActionType, error)
	Get(ctx context.Context, key string) (ActionType, error)
	Delete(ctx context.Context, key string) error
}

type PageUseCase interface {
	Set(ctx context.Context, page Page) error
	GetAll(ctx context.Context) ([]Page, error)
	Get(ctx context.Context, id string) (Page, error)
	Delete(ctx context.Context, id string) error
}

type Dependencies struct {
	ActionUseCase     ActionUseCase
	ActionTypeUseCase ActionTypeUseCase
	PageUseCase       PageUseCase
}
