package mocks

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type PeerStore struct {
	StoreFunc          func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
	UpdatePeerListFunc func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
}

func BaselinePeerStore(t *testing.T) *PeerStore {
	t.Helper()

	peerstore := PeerStore{
		StoreFunc: func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
			return nil
		},
		UpdatePeerListFunc: func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
			return nil
		},
	}

	return &peerstore
}

func (p *PeerStore) Store(id peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {
	return p.StoreFunc(id, addr, info)
}

func (p *PeerStore) UpdatePeerList(id peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {
	return p.StoreFunc(id, addr, info)
}
