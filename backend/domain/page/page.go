package page

import (
	"context"

	"master-control/model"
)

type Repository interface {
	Set(ctx context.Context, page model.Page) error
	GetAll(ctx context.Context) ([]model.Page, error)
	Get(ctx context.Context, ID string) (model.Page, error)
	Delete(ctx context.Context, ID string) error
	Seed(ctx context.Context) error
}

type Page struct {
	repository Repository
}

func New(repository Repository) *Page {
	return &Page{
		repository: repository,
	}
}

func (p *Page) Set(ctx context.Context, page model.Page) error {
	return p.repository.Set(ctx, page)
}

func (p *Page) GetAll(ctx context.Context) ([]model.Page, error) {
	return p.repository.GetAll(ctx)
}

func (p *Page) Get(ctx context.Context, ID string) (model.Page, error) {
	return p.repository.Get(ctx, ID)
}

func (p *Page) Delete(ctx context.Context, ID string) error {
	return p.repository.Delete(ctx, ID)
}