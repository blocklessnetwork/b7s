package store

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// TODO: Implement - RetrievePeers
// TODO: Implement - RetrieveFunctions

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

func (s *Store) RetrievePeer(id peer.ID) (blockless.Peer, error) {

	peerID, err := id.MarshalBinary()
	if err != nil {
		return blockless.Peer{}, fmt.Errorf("could not serialize peer ID: %w", err)
	}

	key := EncodeKey(PrefixPeer, peerID)
	var peer blockless.Peer
	err = s.retrieve(key, &peer)
	if err != nil {
		return blockless.Peer{}, fmt.Errorf("could not retrieve value: %w", err)
	}

	return peer, nil
}

func (s *Store) RetrievePeers() ([]blockless.Peer, error) {

	opts := prefixIterOptions([]byte{PrefixPeer})
	it := s.db.NewIter(opts)

	for it.First(); it.Valid(); it.Next() {

	}
}

// TODO: Define a function record.
func (s *Store) RetrieveFunction(cid string) (any, error) {

	key := EncodeKey(PrefixFunction, cid)

	var function any
	err := s.retrieve(key, &function)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve function record: %w", err)
	}

	return function, nil
}

func (s *Store) RemovePeer(id peer.ID) error {
	return errors.New("TBD: Not implemented")
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

func (s *Store) retrieve(key []byte, out any) error {

	value, closer, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return blockless.ErrNotFound
		}
		return fmt.Errorf("could not retrieve value: %w", err)
	}
	// Closer must be called else a memory leak occurs.
	defer closer.Close()

	err = s.codec.Unmarshal(value, out)
	if err != nil {
		return fmt.Errorf("cold not decode record: %w", err)
	}

	return nil
}

func prefixIterOptions(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: iteratorPrefixUpperBound(prefix),
	}
}

func iteratorPrefixUpperBound(prefix []byte) []byte {

	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}

	return nil
}
