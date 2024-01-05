package node

import (
	"context"
	"encoding/json"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type topicInfo struct {
	handle       *pubsub.Topic
	subscription *pubsub.Subscription
}

func (n *Node) subscribeToTopics(ctx context.Context) error {

	// TODO: If some topics/subscriptions failed, cleanup those already subscribed to.
	for _, topicName := range n.cfg.Topics {

		topic, subscription, err := n.host.Subscribe(ctx, topicName)
		if err != nil {
			return fmt.Errorf("could not subscribe to topic (name: %s): %w", topicName, err)
		}

		ti := &topicInfo{
			handle:       topic,
			subscription: subscription,
		}

		n.topics[topicName] = ti
	}

	return nil
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

// sendToMany serializes the message and sends it to a number of peers. It aborts on any error.
func (n *Node) sendToMany(ctx context.Context, peers []peer.ID, msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	for i, peer := range peers {
		// Send message.
		err = n.host.SendMessage(ctx, peer, payload)
		if err != nil {
			return fmt.Errorf("could not send message to peer (id: %v, peer %d out of %d): %w", peer, i, len(peers), err)
		}
	}

	return nil
}

func (n *Node) publish(ctx context.Context, msg interface{}) error {
	return n.publishToTopic(ctx, DefaultTopic, msg)
}

func (n *Node) publishToTopic(ctx context.Context, topic string, msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	topicInfo, ok := n.topics[topic]
	if !ok {
		return fmt.Errorf("cannot publish to an unknown topic: %s", topic)
	}

	// Publish message.
	err = n.host.Publish(ctx, topicInfo.handle, payload)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	return nil
}

func (n *Node) haveConnection(peer peer.ID) bool {
	connections := n.host.Network().ConnsToPeer(peer)
	return len(connections) > 0
}
