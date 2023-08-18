package peerstore

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// PeerStore takes care of storing and reading peer information to and from persistent storage.
type PeerStore struct {
	store Store
}

// New creates a new PeerStore handler.
func New(store Store) *PeerStore {

	ps := PeerStore{
		store: store,
	}

	return &ps
}

// Get wil retrieve peer with the given ID.
func (p *PeerStore) Get(id peer.ID) (blockless.Peer, error) {

	var peer blockless.Peer
	err := p.store.GetRecord(id.String(), &peer)
	if err != nil {
		return blockless.Peer{}, fmt.Errorf("could not retrieve peer: %w", err)
	}

	return peer, nil
}

// Store will persist the peer information.
func (p *PeerStore) Store(id peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {

	peerInfo := blockless.Peer{
		ID:        id,
		MultiAddr: addr.String(),
		AddrInfo:  info,
	}

	err := p.store.SetRecord(id.String(), peerInfo)
	if err != nil {
		return fmt.Errorf("could not store peer: %w", err)
	}

	return nil
}

// Remove removes the peer from the peerstore.
func (p *PeerStore) Remove(id peer.ID) error {

	err := p.store.Delete(id.String())
	if err != nil {
		return fmt.Errorf("could not remove peer: %w", err)
	}

	return nil
}

// Peers returns the list of peers from the peer store.
func (p *PeerStore) Peers() ([]blockless.Peer, error) {

	ids := p.store.Keys()

	var peers []blockless.Peer
	for _, id := range ids {

		var peer blockless.Peer
		err := p.store.GetRecord(id, &peer)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve peer (id: %v): %w", id, err)
		}

		peers = append(peers, peer)
	}

	return peers, nil
}
