package node

import (
	"fmt"
)

// Tracing span names.
const (
	// message events
	spanMessageSend    = "MessageSend"
	spanMessagePublish = "MessagePublish"
	spanMessageProcess = "MessageProcess"
	// notifiee events
	spanPeerConnected    = "PeerConnected"
	spanPeerDisconnected = "PeerDisconnected"
	// execution events
	spanHeadExecute   = "HeadExecute"
	spanWorkerExecute = "WorkerExecute"
)

// Tracing span status messages.
const (
	spanStatusOK  = "message processed ok"
	spanStatusErr = "error processing message"
)

func msgProcessSpanName(msgType string) string {
	return fmt.Sprintf("%s %s", spanMessageProcess, msgType)
}

func msgSendSpanName(prefix string, msgType string) string {
	return fmt.Sprintf("%s %s", prefix, msgType)
}

var (
	rollCallsSeenMetric        = []string{"node", "rollcalls", "seen"}
	rollCallsPublishedMetric   = []string{"node", "rollcalls", "published"}
	rollCallsAppliedMetric     = []string{"node", "rollcalls", "applied"}
	messagesProcessedMetric    = []string{"node", "messages", "processed"}
	messagesProcessedOkMetric  = []string{"node", "messages", "processed", "ok"}
	messagesProcessedErrMetric = []string{"node", "messages", "processed", "err"}
	functionExecutionsMetric   = []string{"node", "function", "executions"}
	subscriptionsMetric        = []string{"node", "topic", "subscriptions"}
	directMessagesMetric       = []string{"node", "direct", "messages"}
	topicMessagesMetric        = []string{"node", "topic", "messages"}
)
