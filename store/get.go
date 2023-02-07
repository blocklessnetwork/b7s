package store

import (
	"encoding/json"
	"fmt"
)

// Get retrieves the value for a key.
// TODO: Check - do we need both byte and string variants?
// Investigate which ones are more often needed and keep that one.
func (s *Store) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	value, closer, err := s.db.Get([]byte(key))
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

// GetRecord will read and JSON-decode the record associated with the provided key.
func (s *Store) GetRecord(key string, out interface{}) error {
	s.RLock()
	defer s.RUnlock()

	value, closer, err := s.db.Get([]byte(key))
	if err != nil {
		return fmt.Errorf("could not retrieve value: %w", err)
	}
	// Closer must be called else a memory leak occurs.
	defer closer.Close()

	err = json.Unmarshal(value, out)
	if err != nil {
		return fmt.Errorf("could not decode record: %w", err)
	}

	return nil
}
