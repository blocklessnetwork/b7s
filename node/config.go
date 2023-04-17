package node

import (
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	Role:             blockless.WorkerNode,
	Topic:            DefaultTopic,
	HealthInterval:   DefaultHealthInterval,
	RollCallTimeout:  DefaultRollCallTimeout,
	Concurrency:      DefaultConcurrency,
	ExecutionTimeout: DefaultExecutionTimeout,
	Quorum:           DefaultQuorum,
}

// Config represents the Node configuration.
type Config struct {
	Role             blockless.NodeRole // Node role.
	Topic            string             // Topic to subscribe to.
	Execute          Executor           // Executor to use for running functions.
	HealthInterval   time.Duration      // How often should we emit the health ping.
	RollCallTimeout  time.Duration      // How long do we wait for roll call responses.
	Concurrency      uint               // How many requests should the node process in parallel.
	ExecutionTimeout time.Duration      // How long does the head node wait for worker nodes to send their execution results.
	Quorum           uint               // How many nodes do we require for execution.
}

// WithRole specifies the role for the node.
func WithRole(role blockless.NodeRole) Option {
	return func(cfg *Config) {
		cfg.Role = role
	}
}

// WithTopic specifies the p2p topic to which node should subscribe.
func WithTopic(topic string) Option {
	return func(cfg *Config) {
		cfg.Topic = topic
	}
}

// WithExecutor specifies the executor to be used for running Blockless functions
func WithExecutor(execute Executor) Option {
	return func(cfg *Config) {
		cfg.Execute = execute
	}
}

// WithHealthInterval specifies how often we should emit the health signal.
func WithHealthInterval(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.HealthInterval = d
	}
}

// WithRollCallTimeout specifies how long do we wait for roll call responses.
func WithRollCallTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.RollCallTimeout = d
	}
}

// WithExecutionTimeout specifies how long does the head node wait for worker nodes to send their execution results.
func WithExecutionTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.ExecutionTimeout = d
	}
}

// WithQuorum specifies how many worker nodes does the head node want for any given execution.
func WithQuorum(n uint) Option {
	return func(cfg *Config) {
		cfg.Quorum = n
	}
}

// WithConcurrency specifies how many requests the node should process in parallel.
func WithConcurrency(n uint) Option {
	return func(cfg *Config) {
		cfg.Concurrency = n
	}
}
