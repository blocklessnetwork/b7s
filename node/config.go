package node

import (
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	Role:  blockless.WorkerNode,
	Topic: DefaultTopic,
}

// Config represents the Node configuration.
type Config struct {
	Role  blockless.NodeRole // Node role.
	Topic string             // Topic to subscribe to.
}

// WithRole specifies the role for the node.
func WithRole(role blockless.NodeRole) func(*Config) {
	return func(cfg *Config) {
		cfg.Role = role
	}
}

// WithTopic specifies the p2p topic to which node should subscribe.
func WithTopic(topic string) func(*Config) {
	return func(cfg *Config) {
		cfg.Topic = topic
	}
}
