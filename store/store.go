package store

import (
	"github.com/cockroachdb/pebble"
)

// TODO: Consider passing in marshaller/unmarshaller.

// Store enables interaction with a database.
type Store struct {
	db    *pebble.DB
	codec Codec
}

// New creates a new Store backed by the database at the given path.
func New(db *pebble.DB, codec Codec) *Store {

	store := Store{
		db:    db,
		codec: codec,
	}

	return &store
}
