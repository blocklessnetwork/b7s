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

	err := n.host.InitPubSub(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize pubsub: %w", err)
	}

	// TODO: If some topics/subscriptions failed, cleanup those already subscribed to.
	for _, topicName := range n.cfg.Topics {

		topic, subscription, err := n.host.Subscribe(topicName)
		if err != nil {
			return fmt.Errorf("could not subscribe to topic (name: %s): %w", topicName, err)
		}

		ti := &topicInfo{
			handle:       topic,
			subscription: subscription,
		}

		// No need for locking since this initialization is done once on start.
		n.subgroups.topics[topicName] = ti
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

	n.subgroups.RLock()
	topicInfo, ok := n.subgroups.topics[topic]
	n.subgroups.RUnlock()

	if !ok {
		n.log.Info().Str("topic", topic).Msg("unknown topic, joining now")

		var err error
		topicInfo, err = n.joinTopic(topic)
		if err != nil {
			return fmt.Errorf("could not join topic (topic: %s): %w", topic, err)
		}
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
