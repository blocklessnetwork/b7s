package pbft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO (pbft): Consider creating an abstraction like a "network stack" that will be cognizant of the peers in a cluster.
// TODO (pbft): context.Background() used at the moment, fix.

func (r *Replica) send(to peer.ID, msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// Send message.
	err = r.host.SendMessage(context.Background(), to, payload)
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
		err = r.host.SendMessage(context.Background(), peer, payload)
		if err != nil {
			return fmt.Errorf("could not send message: %w", err)
		}
	}

	return nil
}
