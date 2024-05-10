package blockless

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Store interface {
	PeerStore
	FunctionStore
}

type PeerStore interface {
	SavePeer(peer Peer) error
	RetrievePeer(id peer.ID) (Peer, error)
	RetrievePeers() ([]Peer, error)
	RemovePeer(id peer.ID) error
}

type FunctionStore interface {
	SaveFunction(function FunctionRecord) error
	RetrieveFunction(cid string) (FunctionRecord, error)
	RetrieveFunctions() ([]FunctionRecord, error)
	RemoveFunction(id string) error
}
