package node

import (
	"context"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	pp "github.com/blocklessnetwork/b7s/node/internal/pipeline"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

const (
	tracerName = "b7s.Node"
)

func msgProcessSpanOpts(from peer.ID, msgType string, pipeline pp.Pipeline) []trace.SpanStartOption {

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePeer.String(from.String()),
		b7ssemconv.MessageType.String(msgType),
		b7ssemconv.MessagePipeline.String(pipeline.ID.String()),
	}

	if pipeline.ID == pp.PubSub {
		attrs = append(attrs, b7ssemconv.MessageTopic.String(pipeline.Topic))
	}

	return []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	}
}

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

type msgSpanConfig struct {
	msgPipeline pp.Pipeline
	receivers   []peer.ID
}

func (c *msgSpanConfig) pipeline(p pp.Pipeline) *msgSpanConfig {
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

	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePipeline.String(c.msgPipeline.ID.String()),
	}

	if c.msgPipeline.ID == pp.PubSub {
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
