package node

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/blocklessnetworking/b7s/consensus"
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	Role:                    blockless.WorkerNode,
	Topic:                   DefaultTopic,
	HealthInterval:          DefaultHealthInterval,
	RollCallTimeout:         DefaultRollCallTimeout,
	Concurrency:             DefaultConcurrency,
	ExecutionTimeout:        DefaultExecutionTimeout,
	ClusterFormationTimeout: DefaultClusterFormationTimeout,
	DefaultConsensus:        DefaultConsensusAlgorithm,
}

// Config represents the Node configuration.
type Config struct {
	Role                    blockless.NodeRole // Node role.
	Topic                   string             // Topic to subscribe to.
	Execute                 blockless.Executor // Executor to use for running functions.
	HealthInterval          time.Duration      // How often should we emit the health ping.
	RollCallTimeout         time.Duration      // How long do we wait for roll call responses.
	Concurrency             uint               // How many requests should the node process in parallel.
	ExecutionTimeout        time.Duration      // How long does the head node wait for worker nodes to send their execution results.
	ClusterFormationTimeout time.Duration      // How long do we wait for the nodes to form a cluster for an execution.
	Workspace               string             // Directory where we can store files needed for execution.
	DefaultConsensus        consensus.Type     // Default consensus algorithm to use.
}

// Validate checks if the given configuration is correct.
func (n *Node) ValidateConfig() error {

	if !n.cfg.Role.Valid() {
		return errors.New("node role is not valid")
	}

	if n.cfg.Topic == "" {
		return errors.New("topic cannot be empty")
	}

	// Worker specific validation.
	if n.isWorker() {

		if !filepath.IsAbs(n.cfg.Workspace) {
			return errors.New("workspace must be an absolute path")
		}

		// We require an execution component.
		if n.cfg.Execute == nil {
			return errors.New("execution component is required")
		}
	}

	// Head node specific validation.
	if n.isHead() {

		if n.cfg.Execute != nil {
			return errors.New("execution not supported on this type of node")
		}

	}

	return nil
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
func WithExecutor(execute blockless.Executor) Option {
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

// WithClusterFormationTimeout specifies how long does the head node wait for worker nodes to form a consensus cluster.
func WithClusterFormationTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.ClusterFormationTimeout = d
	}
}

// WithConcurrency specifies how many requests the node should process in parallel.
func WithConcurrency(n uint) Option {
	return func(cfg *Config) {
		cfg.Concurrency = n
	}
}

// WithWorkspace specifies the workspace that the node can use for file storage.
func WithWorkspace(path string) Option {
	return func(cfg *Config) {
		cfg.Workspace = path
	}
}

// WithDefaultConsensus specifies the consensus algorithm to use, if not specified in the request.
func WithDefaultConsensus(c consensus.Type) Option {
	return func(cfg *Config) {
		cfg.DefaultConsensus = c
	}
}

func (n *Node) isWorker() bool {
	return n.cfg.Role == blockless.WorkerNode
}

func (n *Node) isHead() bool {
	return n.cfg.Role == blockless.HeadNode
}
