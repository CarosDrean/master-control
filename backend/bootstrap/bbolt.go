package bootstrap

import (
	"master-control/model"

	"go.etcd.io/bbolt"
)

func newBbolt() (*bbolt.DB, error) {
	db, err := bbolt.Open("bbolt.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(model.BucketNameActions))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(model.BucketNameActionTypes))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(model.BucketNamePages))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
