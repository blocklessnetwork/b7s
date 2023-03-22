package store

import (
	"github.com/cockroachdb/pebble"
)

// Store enables interaction with a database.
type Store struct {
	db *pebble.DB
}

// New creates a new Store backed by the database at the given path.
func New(db *pebble.DB) *Store {

	store := Store{
		db: db,
	}

	return &store
}
