package page

import (
	"context"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"

	"master-control/model"
)

var defaults = []model.Page{
	{
		ID:      "general",
		Name:    "General",
		Icon:    "✨",
		Order:   0,
		Columns: 3,
		Rows:    3,
	},
}

type Page struct {
	db         *bbolt.DB
	bucketName string
}

func New(db *bbolt.DB) *Page {
	return &Page{
		db:         db,
		bucketName: model.BucketNamePages,
	}
}

func (p *Page) Set(ctx context.Context, page model.Page) error {
	jsonData, err := json.Marshal(page)
	if err != nil {
		return err
	}

	return p.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(p.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Put([]byte(page.ID), jsonData)
	})
}

func (p *Page) GetAll(ctx context.Context) ([]model.Page, error) {
	pages := []model.Page{}
	err := p.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(p.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var page model.Page
			err := json.Unmarshal(v, &page)
			if err != nil {
				return err
			}

			pages = append(pages, page)
			return nil
		})
	})

	return pages, err
}

func (p *Page) Get(ctx context.Context, ID string) (model.Page, error) {
	var page model.Page
	err := p.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(p.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		v := b.Get([]byte(ID))
		if v == nil {
			return errors.New("page not found")
		}

		return json.Unmarshal(v, &page)
	})

	return page, err
}

func (p *Page) Delete(ctx context.Context, ID string) error {
	return p.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(p.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Delete([]byte(ID))
	})
}

func (p *Page) Seed(ctx context.Context) error {
	return p.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(p.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		for _, d := range defaults {
			data, err := json.Marshal(d)
			if err != nil {
				return err
			}

			if err := b.Put([]byte(d.ID), data); err != nil {
				return err
			}
		}

		return nil
	})
}