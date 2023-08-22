package raft

import (
	"time"
)

// Raft and consensus related parameters.
const (
	defaultConsensusDirName = "consensus"
	defaultLogStoreName     = "logs.dat"
	defaultStableStoreName  = "stable.dat"

	raftClusterDisbandTimeout = 5 * time.Minute
	// Timeout for the context used for sending disband request to cluster nodes.
	raftClusterSendTimeout = 10 * time.Second

	defaultRaftApplyTimeout = 0 // No timeout.
	DefaultHeartbeatTimeout = 300 * time.Millisecond
	DefaultElectionTimeout  = 300 * time.Millisecond
	DefaultLeaderLease      = 200 * time.Millisecond

	consensusTransportTimeout = 1 * time.Minute
)
