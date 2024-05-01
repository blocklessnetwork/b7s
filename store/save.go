package store

import (
	"fmt"

	"github.com/cockroachdb/pebble"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func (s *Store) SavePeer(peer blockless.Peer) error {

	key, err := encodeKey(PrefixPeer, peer.ID)
	if err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	err = s.save(key, peer)
	if err != nil {
		return fmt.Errorf("could not save peer: %w", err)
	}

	return nil
}

func (s *Store) SaveFunction(function blockless.FunctionRecord) error {

	key, err := encodeKey(PrefixFunction, function.CID)
	if err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	err = s.save(key, function)
	if err != nil {
		return fmt.Errorf("could not save function: %w", err)
	}

	return nil
}

func (s *Store) save(key []byte, value any) error {

	encoded, err := s.codec.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not encode value: %w", err)
	}

	err = s.db.Set(key, encoded, pebble.Sync)
	if err != nil {
		return fmt.Errorf("could not store value: %w", err)
	}

	return nil
}
