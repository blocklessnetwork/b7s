package node

import (
	"errors"
	"time"

	"github.com/blocklessnetworking/b7s/consensus"
)

const (
	DefaultTopic                   = "blockless/b7s/general"
	DefaultHealthInterval          = 1 * time.Minute
	DefaultRollCallTimeout         = 5 * time.Second
	DefaultExecutionTimeout        = 10 * time.Second
	DefaultClusterFormationTimeout = 10 * time.Second
	DefaultConcurrency             = 10

	DefaultConsensusAlgorithm = consensus.Raft

	rollCallQueueBufferSize = 1000

	defaultExecutionThreshold = 0.6

	syncInterval = time.Hour // How often do we recheck function installations.
)

// Raft and consensus related parameters.
const (
	defaultConsensusDirName = "consensus"
	defaultLogStoreName     = "logs.dat"
	defaultStableStoreName  = "stable.dat"

	raftClusterDisbandTimeout = 5 * time.Minute
	// Timeout for the context used for sending disband request to cluster nodes.
	raftClusterSendTimeout = 10 * time.Second

	defaultRaftApplyTimeout     = 0 // No timeout.
	DefaultRaftHeartbeatTimeout = 300 * time.Millisecond
	DefaultRaftElectionTimeout  = 300 * time.Millisecond
	DefaultRaftLeaderLease      = 200 * time.Millisecond

	consensusTransportTimeout = 1 * time.Minute
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
)
