package host

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Subscribe will have the host start listening to a specified gossipsub topic.
func (h *Host) Subscribe(ctx context.Context, topic string) (*pubsub.Topic, *pubsub.Subscription, error) {

	// Get a new PubSub object with the default router.
	pubsub, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create new gossipsub: %w", err)
	}

	// Join the specified topic.
	th, err := pubsub.Join(topic)
	if err != nil {
		return nil, nil, fmt.Errorf("could not join topic: %w", err)
	}

	// Subscribe to the topic.
	subscription, err := th.Subscribe()
	if err != nil {
		return nil, nil, fmt.Errorf("could not subscribe to topic: %w", err)
	}

	return th, subscription, nil
}
