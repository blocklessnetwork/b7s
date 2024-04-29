package store

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (s *Store) RemovePeer(id peer.ID) error {

	key := EncodeKey(PrefixPeer, id)
	err := s.remove(key)
	if err != nil {
		return fmt.Errorf("could not remove peer: %w", err)
	}

	return nil
}

func (s *Store) RemoveFunction(cid string) error {

	key := EncodeKey(PrefixFunction, cid)
	err := s.remove(key)
	if err != nil {
		return fmt.Errorf("could not remove function: %w", err)
	}

	return nil
}

func (s *Store) remove(key []byte) error {
	return s.db.Delete(key, pebble.Sync)
}
