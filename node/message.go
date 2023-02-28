package node

import (
	"context"
	"encoding/json"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (n *Node) subscribe(ctx context.Context) (*pubsub.Subscription, error) {

	topic, subscription, err := n.host.Subscribe(ctx, n.topicName)
	if err != nil {
		return nil, fmt.Errorf("could not subscribe to topic: %w", err)
	}
	n.topic = topic

	return subscription, nil
}

// send serializes the message and sends it to the specified peer.
func (n *Node) send(ctx context.Context, to peer.ID, msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// Send message.
	err = n.host.SendMessage(ctx, to, payload)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

func (n *Node) publish(ctx context.Context, msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// Publish message.
	err = n.host.Publish(ctx, n.topic, payload)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	return nil
}
