package pbft

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

const (
	spanAttrView           = attribute.Key("pbft.view")
	spanAttrSequenceNumber = attribute.Key("pbft.sequence_number")
)

func processMessageSpanOptions(from peer.ID, msg any) (string, []trace.SpanStartOption) {

	name := ""
	attrs := []attribute.KeyValue{
		b7ssemconv.MessagePeer.String(from.String()),
	}

	switch m := msg.(type) {

	case Request:
		name = msgSpanName(MessageRequest)
		attrs = append(attrs,
			b7ssemconv.ExecutionRequestID.String(m.ID))

	case PrePrepare:
		name = msgSpanName(MessagePrePrepare)
		attrs = append(attrs,
			spanAttrView.Int64(int64(m.View)),
			spanAttrSequenceNumber.Int64(int64(m.SequenceNumber)),
		)

	case Prepare:
		name = msgSpanName(MessagePrepare)
		attrs = append(attrs,
			spanAttrView.Int64(int64(m.View)),
			spanAttrSequenceNumber.Int64(int64(m.SequenceNumber)))

	case Commit:
		name = msgSpanName(MessageCommit)
		attrs = append(attrs,
			spanAttrView.Int64(int64(m.View)),
			spanAttrSequenceNumber.Int64(int64(m.SequenceNumber)))

	case ViewChange:
		name = msgSpanName(MessageViewChange)
		attrs = append(attrs,
			spanAttrView.Int64(int64(m.View)))

	case NewView:
		name = msgSpanName(MessageNewView)
		attrs = append(attrs,
			spanAttrView.Int64(int64(m.View)))

		// not expecting any other message.
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	}

	return name, opts
}

func msgSpanName(t MessageType) string {
	return fmt.Sprintf("PBFTMessage %s", t.String())
}

func saveTraceContext(ctx context.Context, msg any) {
	tmsg, ok := msg.(blockless.TraceableMessage)
	if ok {
		tmsg.SaveTraceContext(tracing.GetTraceInfo(ctx))
	}
}

// TODO: Unify this and the stuff from the node package and move elsewhere.

// type msgSpanConfig struct {
// 	msgPipeline string
// 	receivers   []peer.ID
// }

// func (c *msgSpanConfig) pipeline(p string) *msgSpanConfig {
// 	c.msgPipeline = p
// 	return c
// }

// func (c *msgSpanConfig) peer(id peer.ID) *msgSpanConfig {
// 	if c.receivers == nil {
// 		c.receivers = make([]peer.ID, 0, 1)
// 	}

// 	c.receivers = append(c.receivers, id)
// 	return c
// }

// func (c *msgSpanConfig) peers(ids ...peer.ID) *msgSpanConfig {
// 	if c.receivers == nil {
// 		c.receivers = make([]peer.ID, 0, len(ids))
// 	}

// 	c.receivers = append(c.receivers, ids...)
// 	return c
// }

// func (c *msgSpanConfig) spanOpts() []trace.SpanStartOption {

// 	var attrs []attribute.KeyValue
// 	if c.msgPipeline != "" {
// 		attrs = append(attrs, b7ssemconv.MessagePipeline.String(c.msgPipeline))
// 	}

// 	if len(c.receivers) == 1 {
// 		attrs = append(attrs, b7ssemconv.MessagePeer.String(c.receivers[0].String()))
// 	} else if len(c.receivers) > 1 {
// 		attrs = append(attrs, b7ssemconv.MessagePeers.String(
// 			strings.Join(blockless.PeerIDsToStr(c.receivers), ","),
// 		))
// 	}

// 	return []trace.SpanStartOption{
// 		trace.WithSpanKind(trace.SpanKindProducer),
// 		trace.WithAttributes(attrs...),
// 	}
// }
