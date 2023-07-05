package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// processMessage will determine which message was received and how to process it.
func (n *Node) processMessage(ctx context.Context, from peer.ID, payload []byte) error {

	// Determine message type.
	msgType, err := getMessageType(payload)
	if err != nil {
		return fmt.Errorf("could not determine message type: %w", err)
	}

	n.log.Trace().Str("peer", from.String()).Str("message", msgType).Msg("received message from peer")

	// Get the registered handler for the message.
	handler := n.getHandler(msgType)

	// Invoke the aprropriate handler to process the message.
	return handler(ctx, from, payload)
}

type baseMessage struct {
	Type string `json:"type,omitempty"`
}

// getMessageType will return the `type` string field from the JSON payload.
func getMessageType(payload []byte) (string, error) {

	var message baseMessage
	err := json.Unmarshal(payload, &message)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal message: %w", err)
	}

	return message.Type, nil
}
