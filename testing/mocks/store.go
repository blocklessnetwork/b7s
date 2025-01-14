package mocks

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/blockless"
)

type Store struct {
	SavePeerFunc      func(context.Context, blockless.Peer) error
	RetrievePeerFunc  func(context.Context, peer.ID) (blockless.Peer, error)
	RetrievePeersFunc func(context.Context) ([]blockless.Peer, error)
	RemovePeerFunc    func(context.Context, peer.ID) error

	SaveFunctionFunc      func(context.Context, blockless.FunctionRecord) error
	RetrieveFunctionFunc  func(context.Context, string) (blockless.FunctionRecord, error)
	RetrieveFunctionsFunc func(context.Context) ([]blockless.FunctionRecord, error)
	RemoveFunctionFunc    func(context.Context, string) error
}

func BaselineStore(t *testing.T) *Store {
	t.Helper()

	store := Store{
		SavePeerFunc: func(context.Context, blockless.Peer) error {
			return nil
		},
		RetrievePeerFunc: func(context.Context, peer.ID) (blockless.Peer, error) {
			return GenericPeer, nil
		},
		RetrievePeersFunc: func(context.Context) ([]blockless.Peer, error) {
			return []blockless.Peer{GenericPeer}, nil
		},
		RemovePeerFunc: func(context.Context, peer.ID) error {
			return nil
		},

		SaveFunctionFunc: func(context.Context, blockless.FunctionRecord) error {
			return nil
		},
		RetrieveFunctionFunc: func(context.Context, string) (blockless.FunctionRecord, error) {
			return GenericFunctionRecord, nil
		},
		RetrieveFunctionsFunc: func(context.Context) ([]blockless.FunctionRecord, error) {
			return []blockless.FunctionRecord{GenericFunctionRecord}, nil
		},
		RemoveFunctionFunc: func(context.Context, string) error {
			return nil
		},
	}

	return &store
}

func (s *Store) SavePeer(ctx context.Context, peer blockless.Peer) error {
	return s.SavePeerFunc(ctx, peer)
}
func (s *Store) SaveFunction(ctx context.Context, function blockless.FunctionRecord) error {
	return s.SaveFunctionFunc(ctx, function)
}
func (s *Store) RetrievePeer(ctx context.Context, id peer.ID) (blockless.Peer, error) {
	return s.RetrievePeerFunc(ctx, id)
}
func (s *Store) RetrievePeers(ctx context.Context) ([]blockless.Peer, error) {
	return s.RetrievePeersFunc(ctx)
}
func (s *Store) RetrieveFunction(ctx context.Context, cid string) (blockless.FunctionRecord, error) {
	return s.RetrieveFunctionFunc(ctx, cid)
}
func (s *Store) RetrieveFunctions(ctx context.Context) ([]blockless.FunctionRecord, error) {
	return s.RetrieveFunctionsFunc(ctx)
}
func (s *Store) RemovePeer(ctx context.Context, id peer.ID) error {
	return s.RemovePeerFunc(ctx, id)
}
func (s *Store) RemoveFunction(ctx context.Context, id string) error {
	return s.RemoveFunctionFunc(ctx, id)
}
