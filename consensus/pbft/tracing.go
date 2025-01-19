package pbft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/telemetry/tracing"
)

const (
	spanMessageProcess   = "MessageProcess"
	spanMessageSend      = "MessageSend"
	spanMessageBroadcast = "MessageBroadcast"
)

func saveTraceContext(ctx context.Context, msg any) {
	tmsg, ok := msg.(bls.TraceableMessage)
	if !ok {
		return
	}

	t := tracing.GetTraceInfo(ctx)
	if !t.Empty() {
		tmsg.SaveTraceContext(t)
	}
}

func msgProcessSpanName(t MessageType) string {
	return fmt.Sprintf("PBFTMessage %s %s", spanMessageProcess, t.String())
}

func msgSendSpanName(msg any, action string) string {
	return fmt.Sprintf("PBFTMessage %s %s", action, messageType(msg))
}

func messageType(msg any) string {
	pmsg, ok := msg.(PBFTMessage)
	if ok {
		return pmsg.Type().String()
	}

	bmsg, ok := msg.(bls.Message)
	if ok {
		return bmsg.Type()
	}

	return ""
}

func getTraceInfoFromMessage(payload []byte) (tracing.TraceInfo, bool) {

	var ti tracing.TraceInfo
	err := json.Unmarshal(payload, &ti)
	if err != nil {
		return ti, false
	}

	// Return true if carrier is populated, false if not.
	return ti, !ti.Empty()
}
