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

	// TODO: (raft) - think abot this
	consensusTransportTimeout = 5 * time.Second

	syncInterval = time.Hour

	// prefix to use for consensus related files and databases.
	consensusDirPrefix = "consensus"

	// TODO: (raft) consider having this configurable
	defaultRaftApplyTimeout = time.Minute

	DefaultRaftHeartbeatTimeout = 300 * time.Millisecond
	DefaultRaftElectionTimeout  = 300 * time.Millisecond
	DefaultRaftLeaderLease      = 200 * time.Millisecond
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
)
