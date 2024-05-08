package main

import (
	"fmt"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func purgeDialbackPeers(store blockless.PeerStore) error {

	// NOTE: If the peer count grows too much in the future - we should adopt a more iterative approach.
	peers, err := store.RetrievePeers()
	if err != nil {
		return fmt.Errorf("could not retrieve peers: %w", err)
	}

	for _, peer := range peers {
		err = store.RemovePeer(peer.ID)
		if err != nil {
			return fmt.Errorf("could not remove peer: %w", err)
		}
	}

	return nil
}
