package node

import (
	"errors"
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

	ClusterAddressTTL = 30 * time.Minute

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
