package host

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Publish will publish the message on the provided gossipsub topic.
func (h *Host) Publish(ctx context.Context, topic *pubsub.Topic, payload []byte) error {

	// Publish the message.
	err := topic.Publish(ctx, payload)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	return nil
}
