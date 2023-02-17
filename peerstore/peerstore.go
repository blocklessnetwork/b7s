package peerstore

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetworking/b7s/models/blockless"
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

// Store will persist the peer information.
func (p *PeerStore) Store(peerID peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {

	// Check if we already have this peer stored.
	var peer blockless.Peer
	err := p.store.GetRecord(peerID.String(), &peer)
	// If we don't have an error it means that the peer is already stored in the DB. We're done.
	if err == nil {
		return nil
	}

	// Check if we failed to retrieve the record. If the error is `not found` - that's okay,
	// and we want to store the peer info now. If it's any other error - halt.
	if err != nil && !errors.Is(err, blockless.ErrNotFound) {
		return fmt.Errorf("could not retrieve peer: %w", err)
	}

	// New peer - create peer info record and store it.
	peerInfo := blockless.Peer{
		Type:      "peer",
		ID:        peerID,
		MultiAddr: addr.String(),
		AddrInfo:  info,
	}

	// Store the peer in the DB.
	// NOTE: This may not be necessary, in case we already had this peer info.
	// Check if we should re-store this data - perhaps we'd be updating part of it?
	err = p.store.SetRecord(peerID.String(), peerInfo)
	if err != nil {
		return fmt.Errorf("could not store peer: %w", err)
	}

	return nil
}

// UpdatePeerList will check if the specified peer is found in the peer list. If not - it will be added.
// NOTE: We're basically duplicating knowledge here - if we have peer stored under its ID,
// we will have it in the `peers` list; do we need to duplicate it?
func (p *PeerStore) UpdatePeerList(peerID peer.ID, addr multiaddr.Multiaddr, info peer.AddrInfo) error {

	// Get list of peers from the store.
	var peers []blockless.Peer
	err := p.store.GetRecord(peersKey, &peers)
	if err != nil && !errors.Is(err, blockless.ErrNotFound) {
		return fmt.Errorf("could not retrieve peer list: %w", err)
	}

	// Check if this is a known peer.
	// NOTE: List iteration, might be slow with a long peer list and frequent connections.
	for _, peer := range peers {

		// If the peer is already known, we're done.
		if peer.ID == peerID {
			return nil
		}
	}

	// New peer - add it to the list of peers.
	peerInfo := blockless.Peer{
		Type:      "peer",
		ID:        peerID,
		MultiAddr: addr.String(),
		AddrInfo:  info,
	}
	peers = append(peers, peerInfo)

	// Store the updated peer list.
	err = p.store.SetRecord(peersKey, peers)
	if err != nil {
		return fmt.Errorf("could not update peer list: %w", err)
	}

	return nil
}

// Peers returns the list of peers from the peer store.
func (p *PeerStore) Peers() ([]blockless.Peer, error) {

	// Get list of peers from the store.
	var peers []blockless.Peer
	err := p.store.GetRecord(peersKey, &peers)
	if err != nil && !errors.Is(err, blockless.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve peer list: %w", err)
	}

	return peers, nil
}
