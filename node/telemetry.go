package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

const (
	tracerName = "b7s.Node"
)

func msgProcessSpanOpts(from peer.ID, msgType string, pipeline messagePipeline) []trace.SpanStartOption {

	// TODO: Topic name - refactor message pipeline to include topic name too.
	// TODO: Message ID is useful but libp2p has dumb IDs.
	// b7ssemconv.MessageID.String(msg.ID),

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			b7ssemconv.MessagePeer.String(from.String()),
			b7ssemconv.MessagePipeline.String(pipeline.String()),
			b7ssemconv.MessageType.String(msgType),
		),
	}
}

// TODO: Use it or lose it.
func traceableTopicName(topic string) string {
	return fmt.Sprintf("topic.%v", topic)
}

func saveTraceContext(ctx context.Context, msg blockless.Message) {
	tmsg, ok := msg.(blockless.TraceableMessage)
	if ok {
		tmsg.SaveTraceContext(tracing.GetTraceInfo(ctx))
	}
}

// TODO: Move this to a separate file.

type msgSpanConfig struct {
	msgPipeline string
	receivers   []peer.ID
}

func (c *msgSpanConfig) pipeline(p string) *msgSpanConfig {
	c.msgPipeline = p
	return c
}

func (c *msgSpanConfig) peer(id peer.ID) *msgSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, 1)
	}

	c.receivers = append(c.receivers, id)
	return c
}

func (c *msgSpanConfig) peers(ids ...peer.ID) *msgSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, len(ids))
	}

	c.receivers = append(c.receivers, ids...)
	return c
}

func (c *msgSpanConfig) spanOpts() []trace.SpanStartOption {

	var attrs []attribute.KeyValue
	if c.msgPipeline != "" {
		attrs = append(attrs, b7ssemconv.MessagePipeline.String(c.msgPipeline))
	}

	if len(c.receivers) == 1 {
		attrs = append(attrs, b7ssemconv.MessagePeer.String(c.receivers[0].String()))
	} else if len(c.receivers) > 1 {
		attrs = append(attrs, b7ssemconv.MessagePeers.String(
			strings.Join(blockless.PeerIDsToStr(c.receivers), ","),
		))
	}

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attrs...),
	}
}
