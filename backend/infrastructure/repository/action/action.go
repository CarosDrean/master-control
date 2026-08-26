package action

import (
	"context"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"

	"master-control/model"
)

type Action struct {
	db         *bbolt.DB
	bucketName string
}

func New(db *bbolt.DB) *Action {
	return &Action{
		db:         db,
		bucketName: model.BucketNameActions,
	}
}

func (a *Action) Set(ctx context.Context, action model.Action) error {
	jsonData, err := json.Marshal(action)
	if err != nil {
		return err
	}

	return a.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Put([]byte(action.ID), jsonData)
	})
}

func (a *Action) GetAll(ctx context.Context) ([]model.Action, error) {
	var actions []model.Action
	err := a.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var action model.Action
			err := json.Unmarshal(v, &action)
			if err != nil {
				return err
			}

			actions = append(actions, action)
			return nil
		})
	})

	return actions, err
}

func (a *Action) Get(ctx context.Context, id string) (model.Action, error) {
	var action model.Action
	err := a.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("action not found")
		}

		return json.Unmarshal(v, &action)
	})

	return action, err
}

func (a *Action) Delete(ctx context.Context, id string) error {
	return a.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Delete([]byte(id))
	})
}
