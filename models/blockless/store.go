package blockless

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Store interface {
	PeerStore
	FunctionStore
}

type PeerStore interface {
	SavePeer(ctx context.Context, peer Peer) error
	RetrievePeer(ctx context.Context, id peer.ID) (Peer, error)
	RetrievePeers(ctx context.Context) ([]Peer, error)
	RemovePeer(ctx context.Context, id peer.ID) error
}

type FunctionStore interface {
	SaveFunction(ctx context.Context, function FunctionRecord) error
	RetrieveFunction(ctx context.Context, cid string) (FunctionRecord, error)
	RetrieveFunctions(ctx context.Context) ([]FunctionRecord, error)
	RemoveFunction(ctx context.Context, id string) error
}
