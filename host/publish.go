package host

import (
	"context"
	"fmt"

	metrics "github.com/armon/go-metrics"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Publish will publish the message on the provided gossipsub topic.
func (h *Host) Publish(ctx context.Context, topic *pubsub.Topic, payload []byte) error {

	metrics.IncrCounterWithLabels(messagesPublishedMetric, 1, []metrics.Label{{Name: "topic", Value: topic.String()}})

	// Publish the message.
	err := topic.Publish(ctx, payload)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	return nil
}
