package node

import (
	"fmt"

	"github.com/armon/go-metrics/prometheus"
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
	rollCallsPublishedMetric   = []string{"node", "rollcalls", "published"}
	rollCallsSeenMetric        = []string{"node", "rollcalls", "seen"}
	rollCallsAppliedMetric     = []string{"node", "rollcalls", "applied"}
	messagesProcessedMetric    = []string{"node", "messages", "processed"}
	messagesProcessedOkMetric  = []string{"node", "messages", "processed", "ok"}
	messagesProcessedErrMetric = []string{"node", "messages", "processed", "err"}
	messagesSentMetric         = []string{"node", "messages", "sent"}
	messagesPublishedMetric    = []string{"node", "messages", "published"}
	functionExecutionsMetric   = []string{"node", "function", "executions"}
	subscriptionsMetric        = []string{"node", "topic", "subscriptions"}
	directMessagesMetric       = []string{"node", "direct", "messages"}
	topicMessagesMetric        = []string{"node", "topic", "messages"}
)

var Counters = []prometheus.CounterDefinition{
	{
		Name: rollCallsPublishedMetric,
		Help: "Number of roll calls this node issued.",
	},
	{
		Name: rollCallsSeenMetric,
		Help: "Number of roll calls seen by the node.",
	},
	{
		Name: rollCallsAppliedMetric,
		Help: "Number of roll calls this node applied to.",
	},
	{
		Name: messagesProcessedMetric,
		Help: "Number of messages this node processed.",
	},
	{
		Name: messagesProcessedOkMetric,
		Help: "Number of messages successfully processed by the node.",
	},
	{
		Name: messagesProcessedErrMetric,
		Help: "Number of messages processed with an error.",
	},
	{
		Name: functionExecutionsMetric,
		Help: "Number of function executions.",
	},
	{
		Name: subscriptionsMetric,
		Help: "Number of topics this node subscribes to.",
	},
	{
		Name: directMessagesMetric,
		Help: "Number of direct messages this node received.",
	},
	{
		Name: topicMessagesMetric,
		Help: "Number of topic messages this node received.",
	},
	{
		Name: messagesSentMetric,
		Help: "Number of messages sent.",
	},
	{
		Name: messagesPublishedMetric,
		Help: "Number of messages published",
	},
}
