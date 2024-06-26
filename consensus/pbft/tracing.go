package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
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
