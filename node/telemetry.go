package node

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

const (
	tracerName = "b7s.Node"
)

func subscriptionMessageSpanOpts(from peer.ID, topicName string) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			// TODO: Message ID is useful but libp2p has dumb IDs.
			// b7ssemconv.MessageID.String(msg.ID),
			b7ssemconv.MessagePeer.String(from.String()),
			b7ssemconv.MessagePipeline.String(traceableTopicName(topicName)),
		),
	}
}

func traceableTopicName(topic string) string {
	return fmt.Sprintf("topic.%v", topic)
}

func directMessageSpanOpts(from peer.ID) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			b7ssemconv.MessagePipeline.String(directMessagePipeline.String()),
			b7ssemconv.MessagePeer.String(from.String()),
		),
	}
}
