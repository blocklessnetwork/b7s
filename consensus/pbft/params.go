package pbft

import (
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	// Protocol to use for PBFT related communication.
	Protocol protocol.ID = "/b7s/consensus/pbft/1.0.0"

	// PBFT offers no resiliency regarding Byzantine nodes with less than four nodes.
	MinimumReplicaCount = 4
)
