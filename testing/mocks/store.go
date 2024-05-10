package mocks

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type Store struct {
	SavePeerFunc          func(blockless.Peer) error
	SaveFunctionFunc      func(blockless.FunctionRecord) error
	RetrievePeerFunc      func(peer.ID) (blockless.Peer, error)
	RetrievePeersFunc     func() ([]blockless.Peer, error)
	RetrieveFunctionFunc  func(string) (blockless.FunctionRecord, error)
	RetrieveFunctionsFunc func() ([]blockless.FunctionRecord, error)
	RemovePeerFunc        func(peer.ID) error
	RemoveFunctionFunc    func(string) error
}

func BaselineStore(t *testing.T) *Store {
	t.Helper()

	store := Store{
		SavePeerFunc: func(blockless.Peer) error {
			return nil
		},
		SaveFunctionFunc: func(blockless.FunctionRecord) error {
			return nil
		},
		RetrievePeerFunc: func(peer.ID) (blockless.Peer, error) {
			return GenericPeer, nil
		},
		RetrievePeersFunc: func() ([]blockless.Peer, error) {
			return []blockless.Peer{GenericPeer}, nil
		},
		RetrieveFunctionFunc: func(string) (blockless.FunctionRecord, error) {
			return GenericFunctionRecord, nil
		},
		RetrieveFunctionsFunc: func() ([]blockless.FunctionRecord, error) {
			return []blockless.FunctionRecord{GenericFunctionRecord}, nil
		},
		RemovePeerFunc: func(peer.ID) error {
			return nil
		},
		RemoveFunctionFunc: func(string) error {
			return nil
		},
	}

	return &store
}

func (s *Store) SavePeer(peer blockless.Peer) error {
	return s.SavePeerFunc(peer)
}
func (s *Store) SaveFunction(function blockless.FunctionRecord) error {
	return s.SaveFunctionFunc(function)
}
func (s *Store) RetrievePeer(id peer.ID) (blockless.Peer, error) {
	return s.RetrievePeerFunc(id)
}
func (s *Store) RetrievePeers() ([]blockless.Peer, error) {
	return s.RetrievePeersFunc()
}
func (s *Store) RetrieveFunction(cid string) (blockless.FunctionRecord, error) {
	return s.RetrieveFunctionFunc(cid)
}
func (s *Store) RetrieveFunctions() ([]blockless.FunctionRecord, error) {
	return s.RetrieveFunctionsFunc()
}
func (s *Store) RemovePeer(id peer.ID) error {
	return s.RemovePeerFunc(id)
}
func (s *Store) RemoveFunction(id string) error {
	return s.RemoveFunctionFunc(id)
}
