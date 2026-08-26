package actiontype

import (
	"context"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"

	"master-control/model"
)

var defaults = []model.ActionType{
	{
		ID:              "app",
		Key:             "app",
		Name:            "App",
		Description:     "Launch a macOS application",
		CommandTemplate: "open -a",
	},
	{
		ID:              "command",
		Key:             "command",
		Name:            "Command",
		Description:     "Execute a shell command",
		CommandTemplate: "bash -c",
	},
	{
		ID:              "applescript",
		Key:             "applescript",
		Name:            "AppleScript",
		Description:     "Execute an inline AppleScript",
		CommandTemplate: "osascript -e",
	},
	{
		ID:              "hotkey",
		Key:             "hotkey",
		Name:            "Hotkey",
		Description:     "Send keyboard shortcuts via System Events",
		CommandTemplate: "osascript -e",
	},
	{
		ID:              "system",
		Key:             "system",
		Name:            "System",
		Description:     "Control volume, brightness and media natively",
		CommandTemplate: "osascript -e",
	},
}

type ActionType struct {
	db         *bbolt.DB
	bucketName string
}

func New(db *bbolt.DB) *ActionType {
	return &ActionType{
		db:         db,
		bucketName: model.BucketNameActionTypes,
	}
}

func (a *ActionType) Set(ctx context.Context, actionType model.ActionType) error {
	jsonData, err := json.Marshal(actionType)
	if err != nil {
		return err
	}

	return a.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Put([]byte(actionType.ID), jsonData)
	})
}

func (a *ActionType) GetAll(ctx context.Context) ([]model.ActionType, error) {
	var actionTypes []model.ActionType
	err := a.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var at model.ActionType
			err := json.Unmarshal(v, &at)
			if err != nil {
				return err
			}

			actionTypes = append(actionTypes, at)
			return nil
		})
	})

	return actionTypes, err
}

func (a *ActionType) Get(ctx context.Context, ID string) (model.ActionType, error) {
	var actionType model.ActionType
	err := a.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		v := b.Get([]byte(ID))
		if v == nil {
			return errors.New("action type not found")
		}

		return json.Unmarshal(v, &actionType)
	})

	return actionType, err
}

func (a *ActionType) Delete(ctx context.Context, ID string) error {
	return a.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		return b.Delete([]byte(ID))
	})
}

func (a *ActionType) Seed(ctx context.Context) error {
	return a.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(a.bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		valid := make(map[string]struct{}, len(defaults))
		for _, d := range defaults {
			valid[d.ID] = struct{}{}

			data, err := json.Marshal(d)
			if err != nil {
				return err
			}

			if err := b.Put([]byte(d.ID), data); err != nil {
				return err
			}
		}

		// remove stale action types no longer defined as defaults
		stale := [][]byte{}
		if err := b.ForEach(func(k, v []byte) error {
			if _, ok := valid[string(k)]; !ok {
				stale = append(stale, k)
			}
			return nil
		}); err != nil {
			return err
		}

		for _, k := range stale {
			if err := b.Delete(k); err != nil {
				return err
			}
		}

		return nil
	})
}
