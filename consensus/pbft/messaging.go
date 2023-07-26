package pbft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// TODO (pbft): Consider creating an abstraction like a "network stack" that will be cognizant of the peers in a cluster.
// TODO (pbft): context.Background() used at the moment, fix.

func (r *Replica) send(to peer.ID, msg interface{}, protocol protocol.ID) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// Send message.
	err = r.host.SendMessageOnProtocol(context.Background(), to, payload, protocol)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

// broadcast sends message to all peers in the replica set.
func (r *Replica) broadcast(msg interface{}) error {

	// TODO (pbft): Harden this. Consider what to do when some sends fail. Say we failed or retry?
	//	It's a valid scenario that some peers may be offline, so we should probably continue
	//	despite errors, at least to a point.

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	for _, peer := range r.peers {

		// Skip self.
		if peer == r.id {
			continue
		}

		err = r.host.SendMessageOnProtocol(context.Background(), peer, payload, Protocol)
		if err != nil {
			return fmt.Errorf("could not send message: %w", err)
		}
	}

	return nil
}
