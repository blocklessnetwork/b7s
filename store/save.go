package store

import (
	"fmt"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/cockroachdb/pebble"
)

func (s *Store) SavePeer(peer blockless.Peer) error {

	id, err := peer.ID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("could not serialize peer ID: %w", err)
	}

	key := EncodeKey(PrefixPeer, id)
	err = s.save(key, peer)
	if err != nil {
		return fmt.Errorf("could not save peer: %w", err)
	}

	return nil
}

// TODO: Define a function record
func (s *Store) SaveFunction(cid string, record any) error {

	key := EncodeKey(PrefixFunction, cid)
	err := s.save(key, record)
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
