package store

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/pebble"
)

// Set sets the value for a key.
func (s *Store) Set(key string, value string) error {

	err := s.db.Set([]byte(key), []byte(value), pebble.Sync)
	if err != nil {
		return fmt.Errorf("could not store value: %w", err)
	}

	return nil
}

// SetRecord will JSON-encode the provided record and store it in the DB.
func (s *Store) SetRecord(key string, value interface{}) error {

	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not serialize the record: %w", err)
	}

	err = s.db.Set([]byte(key), encoded, pebble.Sync)
	if err != nil {
		return fmt.Errorf("could not store value: %w", err)
	}

	return nil
}
