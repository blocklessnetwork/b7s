package mocks

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type PeerStore struct {
	StoreFunc func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
	PeersFunc func() ([]blockless.Peer, error)
}

func BaselinePeerStore(t *testing.T) *PeerStore {
	t.Helper()

	peerstore := PeerStore{
		StoreFunc: func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
			return nil
		},

		PeersFunc: func() ([]blockless.Peer, error) {
			return []blockless.Peer{}, nil
		},
	}

	return &peerstore
}

func (p *PeerStore) Store(id peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {
	return p.StoreFunc(id, addr, info)
}

func (p *PeerStore) Peers() ([]blockless.Peer, error) {
	return p.PeersFunc()
}
