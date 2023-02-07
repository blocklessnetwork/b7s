package database

import (
	"fmt"
	"sync"

	"github.com/cockroachdb/pebble"
)

// TODO: Perhaps name it `storage` instead of `database`.
// TODO: Check - do we need a RWMutex for DB access?

// DB enables interaction with a database.
type DB struct {
	sync.RWMutex

	db *pebble.DB
}

// Connect establishes a connection to a database at the given path.
func Connect(path string) (*DB, error) {

	opts := pebble.Options{}
	pdb, err := pebble.Open(path, &opts)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	db := DB{
		db: pdb,
	}

	return &db, nil
}

// Set sets the value for a key.
func (db *DB) Set(key string, value string) error {
	db.Lock()
	defer db.Unlock()

	err := db.db.Set([]byte(key), []byte(value), pebble.Sync)
	if err != nil {
		return fmt.Errorf("could not save value: %w", err)
	}

	return nil
}

// Get retrieves the value for a key.
// TODO: Check - do we need both byte and string variants?
// Investigate which ones are more often needed and keep that one.
func (db *DB) Get(key string) (string, error) {
	db.RLock()
	defer db.RUnlock()

	value, closer, err := db.db.Get([]byte(key))
	if err != nil {
		return "", fmt.Errorf("could not retrieve value: %w", err)
	}
	// Closer must be called else a memory leak occurs.
	defer closer.Close()

	// After closer is done, the slice is no longer valid, so we need to copy it.
	dup := make([]byte, len(value))
	copy(dup, value)

	return string(dup), nil
}

// Close closes the DB connection.
func (db *DB) Close() error {
	db.Lock()
	defer db.Unlock()

	err := db.db.Close()
	if err != nil {
		return fmt.Errorf("could not close the DB: %w", err)
	}

	return nil
}
