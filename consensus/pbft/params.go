package pbft

import (
	"errors"
	"time"

	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	// Protocol to use for PBFT related communication.
	Protocol protocol.ID = "/b7s/consensus/pbft/1.0.0"

	// PBFT offers no resiliency towards Byzantine nodes with less than four nodes.
	MinimumReplicaCount = 4

	// How long do the send/broadcast operation have until we consider it failed.
	// TODO: Check - doesn't this go against the claim in PBFT that messages eventually get delivered? Think how to handle this.
	NetworkTimeout = 5 * time.Second

	// How long is the inactivity period before we trigger a view change.
	RequestTimeout = 10 * time.Second
)

var (
	ErrViewChange = errors.New("view change in progress")
)

var (
	NullRequest = Request{}
)
