package store

import (
	"sync"

	"github.com/cockroachdb/pebble"
)

// TODO: Check - do we need a RWMutex for DB access?

// Store enables interaction with a database.
type Store struct {
	sync.RWMutex

	db *pebble.DB
}

// New creates a new Store backed by the database at the given path.
func New(db *pebble.DB) (*Store, error) {

	store := Store{
		db: db,
	}

	return &store, nil
}
