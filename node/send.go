package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// send serializes the response and sends it to the specified peer.
func (n *Node) send(ctx context.Context, to peer.ID, res interface{}) error {

	// Serialize the response.
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not encode the record: %w", err)
	}

	// Send message.
	err = n.host.SendMessage(ctx, to, payload)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}
