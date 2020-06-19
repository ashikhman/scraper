package entity

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

const (
	sourceUuidPrefix = "source.uuid."
)

type source struct {
	id     uuid.UUID
	Domain string
}

func NewSource() *source {
	return &source{
		id: uuid.New(),
	}
}

type sourceRepository struct {
	db *badger.DB
}

func NewSourceRepository(db *badger.DB) *sourceRepository {
	return &sourceRepository{
		db: db,
	}
}

func (r *sourceRepository) Save(source *source) error {
	return r.db.Update(func(tx *badger.Txn) error {
		var key = key(sourceUuidPrefix, source.id.String())

		item, err := tx.Get(key)
		if err == nil {
			err = existingItem.Value(func(existing []byte) error {
				return decode(existing, existingVal)
			})
			if err != nil {
				return err
			}
		} else if err != badger.ErrKeyNotFound {
			return err
		}
	})
}

func (r *sourceRepository) Get(id uuid.UUID) (source source, err error) {
	err = r.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key(sourceUuidPrefix, id.String()))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return decode(val, &source)
		})
	})

	return
}
