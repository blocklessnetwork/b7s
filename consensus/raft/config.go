package raft

import (
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/log/hclog"
)

// Option can be used to set Raft configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the raft handler.
var DefaultConfig = Config{
	HeartbeatTimeout: DefaultHeartbeatTimeout,
	ElectionTimeout:  DefaultElectionTimeout,
	LeaderLease:      DefaultLeaderLease,
}

type Config struct {
	Callbacks []FSMProcessFunc // Callback functions to be invoked by the FSM after execution is done.

	HeartbeatTimeout time.Duration // How often a consensus cluster leader should ping its followers.
	ElectionTimeout  time.Duration // How long does a consensus cluster node wait for a leader before it triggers an election.
	LeaderLease      time.Duration // How long does a leader remain a leader if it cannot contact a quorum of cluster nodes.
}

// WithHeartbeatTimeout sets the heartbeat timeout for the consensus cluster.
func WithHeartbeatTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.HeartbeatTimeout = d
	}
}

// WithElectionTimeout sets the election timeout for the consensus cluster.
func WithElectionTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.ElectionTimeout = d
	}
}

// WithLeaderLease sets the leader lease for the consensus cluster leader.
func WithLeaderLease(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.LeaderLease = d
	}
}

func WithCallbacks(callbacks ...FSMProcessFunc) Option {
	return func(cfg *Config) {
		var fns []FSMProcessFunc
		fns = append(fns, callbacks...)
		cfg.Callbacks = fns
	}
}

func getRaftConfig(cfg Config, log zerolog.Logger, nodeID string) raft.Config {

	rcfg := raft.DefaultConfig()
	rcfg.LocalID = raft.ServerID(nodeID)
	rcfg.Logger = hclog.New(log).Named("raft")
	rcfg.HeartbeatTimeout = cfg.HeartbeatTimeout
	rcfg.ElectionTimeout = cfg.ElectionTimeout
	rcfg.LeaderLeaseTimeout = cfg.LeaderLease

	return *rcfg
}

func consensusDir(workspace string, requestID string) string {
	return filepath.Join(workspace, defaultConsensusDirName, requestID)
}
