package node

import (
	"context"
	"fmt"
)

// TODO: Check import - is it `go-libp2p-pubsub` (analogous to `go-libp2p-core`)?

// Run will start the main loop for the node.
func (n Node) Run(ctx context.Context) error {

	// Subscribe to the specified topic.
	topic, subscription, err := n.host.Subscribe(ctx, n.topicName)
	if err != nil {
		return fmt.Errorf("could not subscribe to topic: %w", err)
	}

	// TODO: Perhaps have Host handle this, and just use topic by name?
	n.topic = topic

	// TODO: Stop condition.

	// Process messages.
	for {
		// Receive message.
		msg, err := subscription.Next(ctx)
		if err != nil {
			n.log.Error().
				Err(err).
				Msg("could not receive message")
			continue
		}

		// Skip messages we published.
		if msg.ReceivedFrom == n.host.ID() {
			// TODO: Check there's a field msg.Local - is that the same as ID comparison?
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
