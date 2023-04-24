package mocks

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type PeerStore struct {
	GetFunc    func(peer.ID) (blockless.Peer, error)
	StoreFunc  func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
	PeersFunc  func() ([]blockless.Peer, error)
	RemoveFunc func(peer.ID) error
}

func BaselinePeerStore(t *testing.T) *PeerStore {
	t.Helper()

	peerstore := PeerStore{
		GetFunc: func(peer.ID) (blockless.Peer, error) {
			return blockless.Peer{}, nil
		},
		StoreFunc: func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
			return nil
		},
		PeersFunc: func() ([]blockless.Peer, error) {
			return []blockless.Peer{}, nil
		},
		RemoveFunc: func(peer.ID) error {
			return GenericError
		},
	}

	return &peerstore
}

func (p *PeerStore) Get(id peer.ID) (blockless.Peer, error) {
	return p.GetFunc(id)
}

func (p *PeerStore) Store(id peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {
	return p.StoreFunc(id, addr, info)
}

func (p *PeerStore) Peers() ([]blockless.Peer, error) {
	return p.PeersFunc()
}

func (p *PeerStore) Remove(id peer.ID) error {
	return p.RemoveFunc(id)
}
