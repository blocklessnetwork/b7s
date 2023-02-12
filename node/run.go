package node

import (
	"context"
	"fmt"
)

// Run will start the main loop for the node.
func (n Node) Run(ctx context.Context) error {

	// Subscribe to the specified topic.
	topic, subscription, err := n.host.Subscribe(ctx, n.topicName)
	if err != nil {
		return fmt.Errorf("could not subscribe to topic: %w", err)
	}
	n.topic = topic

	// Message processing loop.
	for {

		// Retrieve next message.
		msg, err := subscription.Next(ctx)
		if err != nil {
			// NOTE: Cancelling the context will lead us here.
			n.log.Error().Err(err).Msg("could not receive message")
			break
		}

		// Skip messages we published.
		if msg.ReceivedFrom == n.host.ID() {
			continue
		}

		n.log.Debug().
			Str("id", msg.ID).
			Str("peer_id", msg.ReceivedFrom.String()).
			Msg("received message")

		err = n.processMessage(ctx, msg.ReceivedFrom, msg.Data)
		if err != nil {
			n.log.Error().
				Err(err).
				Str("id", msg.ID).
				Str("peer_id", msg.ReceivedFrom.String()).
				Msg("could not process message")
			continue
		}
	}

	return nil
}
