package host

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (h *Host) InitPubSub(ctx context.Context) error {

	// Get a new PubSub object with the default router.
	pubsub, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return fmt.Errorf("could not create new gossipsub: %w", err)
	}
	h.pubsub = pubsub

	return nil
}

// Subscribe will have the host start listening to a specified gossipsub topic.
func (h *Host) Subscribe(topic string) (*pubsub.Topic, *pubsub.Subscription, error) {

	// Join the specified topic.
	th, err := h.pubsub.Join(topic)
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
