package store

import (
	"fmt"

	"github.com/cockroachdb/pebble"
)

// Delete removes the key from the database.
func (s *Store) Delete(key string) error {

	err := s.db.Delete([]byte(key), pebble.Sync)
	if err != nil {
		return fmt.Errorf("could not delete value: %w", err)
	}

	return nil
}
