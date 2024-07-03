package node

import (
	"errors"
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/consensus"
)

const (
	DefaultTopic                   = "blockless/b7s/general"
	DefaultHealthInterval          = 1 * time.Minute
	DefaultRollCallTimeout         = 5 * time.Second
	DefaultExecutionTimeout        = 20 * time.Second
	DefaultClusterFormationTimeout = 10 * time.Second
	DefaultConcurrency             = 10

	DefaultConsensusAlgorithm = consensus.Raft

	DefaultAttributeLoadingSetting = false

	rollCallQueueBufferSize = 1000

	defaultExecutionThreshold = 0.6

	syncInterval = time.Hour // How often do we recheck function installations.

	allowErrorLeakToTelemetry = false // By default we will not send processing errors to telemetry tracers.
)

// Raft and consensus related parameters.
const (
	// When disbanding a cluster, how long do we wait until a potential execution is done.
	consensusClusterDisbandTimeout = 5 * time.Minute
	// Timeout for the context used for sending disband request to cluster nodes.
	consensusClusterSendTimeout = 10 * time.Second
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
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
