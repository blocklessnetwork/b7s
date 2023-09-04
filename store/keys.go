package store

import (
	"fmt"
)

// Keys returns the list of all keys in the database.
func (s *Store) Keys() ([]string, error) {

	it, err := s.db.NewIter(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new iterator: %w", err)
	}

	var keys []string
	for it.First(); it.Valid(); it.Next() {

		key := string(it.Key())
		keys = append(keys, key)
	}

	return keys, nil
}
