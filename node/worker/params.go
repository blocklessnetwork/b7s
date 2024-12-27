package worker

import (
	"time"
)

const (
	DefaultAttributeLoadingSetting = false

	ClusterAddressTTL = 30 * time.Minute

	consensusClusterSendTimeout = 10 * time.Second

	syncInterval = time.Hour // How often do we recheck function installations.
)

// Raft and consensus related parameters.
const (
	// When disbanding a cluster, how long do we wait until a potential execution is done.
	consensusClusterDisbandTimeout = 5 * time.Minute
)
