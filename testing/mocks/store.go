package mocks

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/bls"
)

type Store struct {
	SavePeerFunc      func(context.Context, bls.Peer) error
	RetrievePeerFunc  func(context.Context, peer.ID) (bls.Peer, error)
	RetrievePeersFunc func(context.Context) ([]bls.Peer, error)
	RemovePeerFunc    func(context.Context, peer.ID) error

	SaveFunctionFunc      func(context.Context, bls.FunctionRecord) error
	RetrieveFunctionFunc  func(context.Context, string) (bls.FunctionRecord, error)
	RetrieveFunctionsFunc func(context.Context) ([]bls.FunctionRecord, error)
	RemoveFunctionFunc    func(context.Context, string) error
}

func BaselineStore(t *testing.T) *Store {
	t.Helper()

	store := Store{
		SavePeerFunc: func(context.Context, bls.Peer) error {
			return nil
		},
		RetrievePeerFunc: func(context.Context, peer.ID) (bls.Peer, error) {
			return GenericPeer, nil
		},
		RetrievePeersFunc: func(context.Context) ([]bls.Peer, error) {
			return []bls.Peer{GenericPeer}, nil
		},
		RemovePeerFunc: func(context.Context, peer.ID) error {
			return nil
		},

		SaveFunctionFunc: func(context.Context, bls.FunctionRecord) error {
			return nil
		},
		RetrieveFunctionFunc: func(context.Context, string) (bls.FunctionRecord, error) {
			return GenericFunctionRecord, nil
		},
		RetrieveFunctionsFunc: func(context.Context) ([]bls.FunctionRecord, error) {
			return []bls.FunctionRecord{GenericFunctionRecord}, nil
		},
		RemoveFunctionFunc: func(context.Context, string) error {
			return nil
		},
	}

	return &store
}

func (s *Store) SavePeer(ctx context.Context, peer bls.Peer) error {
	return s.SavePeerFunc(ctx, peer)
}
func (s *Store) SaveFunction(ctx context.Context, function bls.FunctionRecord) error {
	return s.SaveFunctionFunc(ctx, function)
}
func (s *Store) RetrievePeer(ctx context.Context, id peer.ID) (bls.Peer, error) {
	return s.RetrievePeerFunc(ctx, id)
}
func (s *Store) RetrievePeers(ctx context.Context) ([]bls.Peer, error) {
	return s.RetrievePeersFunc(ctx)
}
func (s *Store) RetrieveFunction(ctx context.Context, cid string) (bls.FunctionRecord, error) {
	return s.RetrieveFunctionFunc(ctx, cid)
}
func (s *Store) RetrieveFunctions(ctx context.Context) ([]bls.FunctionRecord, error) {
	return s.RetrieveFunctionsFunc(ctx)
}
func (s *Store) RemovePeer(ctx context.Context, id peer.ID) error {
	return s.RemovePeerFunc(ctx, id)
}
func (s *Store) RemoveFunction(ctx context.Context, id string) error {
	return s.RemoveFunctionFunc(ctx, id)
}
