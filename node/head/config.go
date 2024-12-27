package head

import (
	"time"

	"github.com/blocklessnetwork/b7s/consensus"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	RollCallTimeout:         DefaultRollCallTimeout,
	ExecutionTimeout:        DefaultExecutionTimeout,
	ClusterFormationTimeout: DefaultClusterFormationTimeout,
	DefaultConsensus:        DefaultConsensusAlgorithm,
}

// Config represents the Node configuration.
type Config struct {
	RollCallTimeout         time.Duration  // How long do we wait for roll call responses.
	ExecutionTimeout        time.Duration  // How long does the head node wait for worker nodes to send their execution results.
	ClusterFormationTimeout time.Duration  // How long do we wait for the nodes to form a cluster for an execution.
	DefaultConsensus        consensus.Type // Default consensus algorithm to use.
}

func (c Config) Valid() error {
	return nil
}
