package node

import (
	"context"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

func saveTraceContext(ctx context.Context, msg blockless.Message) {
	tmsg, ok := msg.(blockless.TraceableMessage)
	if !ok {
		return
	}

	t := tracing.GetTraceInfo(ctx)
	if !t.Empty() {
		tmsg.SaveTraceContext(t)
	}
}

type messageSpanConfig struct {
	msgPipeline Pipeline
	receivers   []peer.ID
}

func (c *messageSpanConfig) pipeline(p Pipeline) *messageSpanConfig {
	c.msgPipeline = p
	return c
}

func (c *messageSpanConfig) peer(id peer.ID) *messageSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, 1)
	}

	c.receivers = append(c.receivers, id)
	return c
}

func (c *messageSpanConfig) peers(ids ...peer.ID) *messageSpanConfig {
	if c.receivers == nil {
		c.receivers = make([]peer.ID, 0, len(ids))
	}

	c.receivers = append(c.receivers, ids...)
	return c
}

func (c *messageSpanConfig) spanOpts() []trace.SpanStartOption {

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePipeline.String(c.msgPipeline.ID.String()),
	}

	if c.msgPipeline.ID == PubSub {
		attrs = append(attrs, b7ssemconv.MessageTopic.String(c.msgPipeline.Topic))
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

func msgProcessSpanOpts(from peer.ID, msgType string, pipeline Pipeline) []trace.SpanStartOption {

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePeer.String(from.String()),
		b7ssemconv.MessageType.String(msgType),
		b7ssemconv.MessagePipeline.String(pipeline.ID.String()),
	}

	if pipeline.ID == PubSub {
		attrs = append(attrs, b7ssemconv.MessageTopic.String(pipeline.Topic))
	}

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	}
}
