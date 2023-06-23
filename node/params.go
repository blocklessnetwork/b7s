package node

import (
	"errors"
	"time"
)

const (
	DefaultTopic                   = "blockless/b7s/general"
	DefaultHealthInterval          = 1 * time.Minute
	DefaultRollCallTimeout         = 5 * time.Second
	DefaultExecutionTimeout        = 10 * time.Second
	DefaultClusterFormationTimeout = 10 * time.Second
	DefaultConcurrency             = 10

	rollCallQueueBufferSize = 1000

	syncInterval = time.Hour // How often do we recheck function installations.
)

// Raft and consensus related parameters.
const (
	defaultConsensusDirName = "consensus"
	defaultLogStoreName     = "logs.dat"
	defaultStableStoreName  = "stable.dat"

	defaultRaftApplyTimeout     = 0 // No timeout.
	DefaultRaftHeartbeatTimeout = 300 * time.Millisecond
	DefaultRaftElectionTimeout  = 300 * time.Millisecond
	DefaultRaftLeaderLease      = 200 * time.Millisecond

	consensusTransportTimeout = 1 * time.Minute
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
)
