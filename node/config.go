package node

import (
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	Role:            blockless.WorkerNode,
	Topic:           DefaultTopic,
	HealthInterval:  DefaultHealthInterval,
	RollCallTimeout: DefaultRollCallTimeout,
}

// Config represents the Node configuration.
type Config struct {
	Role            blockless.NodeRole // Node role.
	Topic           string             // Topic to subscribe to.
	Execute         Executor           // Executor to use for running functions.
	HealthInterval  time.Duration      // How often should we emit the health ping.
	RollCallTimeout time.Duration      // How long do we wait for roll call responses.
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
