package store

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func (s *Store) RetrievePeer(id peer.ID) (blockless.Peer, error) {

	key, err := encodeKey(PrefixPeer, id)
	if err != nil {
		return blockless.Peer{}, fmt.Errorf("could not encode key: %w", err)
	}
	var peer blockless.Peer
	err = s.retrieve(key, &peer)
	if err != nil {
		return blockless.Peer{}, fmt.Errorf("could not retrieve value: %w", err)
	}

	return peer, nil
}

func (s *Store) RetrievePeers() ([]blockless.Peer, error) {

	peers := make([]blockless.Peer, 0)

	opts := prefixIterOptions([]byte{PrefixPeer})
	it := s.db.NewIter(opts)
	for it.First(); it.Valid(); it.Next() {

		var peer blockless.Peer
		err := s.retrieve(it.Key(), &peer)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve peer (key: %x): %w", it.Key(), err)
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

func (s *Store) RetrieveFunction(cid string) (blockless.FunctionRecord, error) {

	key, err := encodeKey(PrefixFunction, cid)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not encode key: %w", err)
	}

	var function blockless.FunctionRecord
	err = s.retrieve(key, &function)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not retrieve function record: %w", err)
	}

	return function, nil
}

func (s *Store) RetrieveFunctions() ([]blockless.FunctionRecord, error) {

	functions := make([]blockless.FunctionRecord, 0)

	opts := prefixIterOptions([]byte{PrefixFunction})
	it := s.db.NewIter(opts)
	for it.First(); it.Valid(); it.Next() {

		var function blockless.FunctionRecord
		err := s.retrieve(it.Key(), &function)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve functioN (key: %x): %w", it.Key(), err)
		}

		functions = append(functions, function)
	}

	return functions, nil
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
