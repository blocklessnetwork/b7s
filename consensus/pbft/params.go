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
	NetworkTimeout = 5 * time.Second

	// How long is the inactivity period before we trigger a view change.
	RequestTimeout = 10 * time.Second

	EnvVarByzantine = "B7S_PBFT_BYZANTINE"
)

var (
	ErrViewChange            = errors.New("view change in progress")
	ErrActiveView            = errors.New("replica is currently in an active view")
	ErrConflictingPreprepare = errors.New("conflicting pre-prepare")
	ErrInvalidSignature      = errors.New("invalid signature")
)

var (
	NullRequest = Request{}
)
