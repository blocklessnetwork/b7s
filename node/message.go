package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-multierror"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
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

	n.log.Info().Strs("topics", n.cfg.Topics).Msg("topics node will subscribe to")

	metrics.IncrCounter(subscriptionsMetric, float32(len(n.cfg.Topics)))

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
func (n *Node) send(ctx context.Context, to peer.ID, msg blockless.Message) error {

	opts := new(msgSpanConfig).pipeline(directMessagePipeline.String()).peer(to).spanOpts()
	ctx, span := n.tracer.Start(ctx, msgSendSpanName(spanMessageSend, msg.Type()), opts...)
	defer span.End()

	saveTraceContext(ctx, msg)

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

	metrics.IncrCounterWithLabels(messagesSentMetric, 1, []metrics.Label{{Name: "type", Value: msg.Type()}})

	return nil
}

// sendToMany serializes the message and sends it to a number of peers. `requireAll` dictates how we treat partial errors.
func (n *Node) sendToMany(ctx context.Context, peers []peer.ID, msg blockless.Message, requireAll bool) error {

	opts := new(msgSpanConfig).pipeline(directMessagePipeline.String()).peers(peers...).spanOpts()
	ctx, span := n.tracer.Start(ctx, msgSendSpanName(spanMessageSend, msg.Type()), opts...)
	defer span.End()

	saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	var errGroup multierror.Group
	for i, peer := range peers {
		i := i
		peer := peer

		errGroup.Go(func() error {
			err := n.host.SendMessage(ctx, peer, payload)
			if err != nil {
				return fmt.Errorf("peer %v/%v send error (peer: %v): %w", i+1, len(peers), peer.String(), err)
			}

			return nil
		})
	}

	metrics.IncrCounterWithLabels(messagesSentMetric, float32(len(peers)), []metrics.Label{{Name: "type", Value: msg.Type()}})

	retErr := errGroup.Wait()
	if retErr == nil || len(retErr.Errors) == 0 {
		// If everything succeeded => ok.
		return nil
	}

	switch len(retErr.Errors) {
	case len(peers):
		// If everything failed => error.
		return fmt.Errorf("all sends failed: %w", retErr)

	default:
		// Some sends failed - do as requested by `requireAll`.
		if requireAll {
			return fmt.Errorf("some sends failed: %w", retErr)
		}

		n.log.Warn().Err(retErr).Msg("some sends failed, proceeding")

		return nil
	}
}

func (n *Node) publish(ctx context.Context, msg blockless.Message) error {
	return n.publishToTopic(ctx, DefaultTopic, msg)
}

func (n *Node) publishToTopic(ctx context.Context, topic string, msg blockless.Message) error {

	opts := new(msgSpanConfig).pipeline(traceableTopicName(topic)).spanOpts()
	ctx, span := n.tracer.Start(ctx, msgSendSpanName(spanMessagePublish, msg.Type()), opts...)
	defer span.End()

	saveTraceContext(ctx, msg)

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

	metrics.IncrCounterWithLabels(messagesPublishedMetric, 1,
		[]metrics.Label{
			{Name: "type", Value: msg.Type()},
			{Name: "topic", Value: topic},
		})

	return nil
}

func (n *Node) haveConnection(peer peer.ID) bool {
	connections := n.host.Network().ConnsToPeer(peer)
	return len(connections) > 0
}
