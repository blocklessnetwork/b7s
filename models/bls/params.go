package bls

import (
	"errors"
	"time"

	"github.com/libp2p/go-libp2p/core/protocol"
)

// Sentinel errors.
var (
	ErrNotFound                = errors.New("not found")
	ErrRollCallTimeout         = errors.New("roll call timed out - not enough nodes responded")
	ErrExecutionNotEnoughNodes = errors.New("not enough execution results received")
)

const (
	ProtocolID protocol.ID = "/b7s/work/1.0.0"
	EnvPrefix  string      = "B7S_"

	DefaultTopic          = "blockless/b7s/general"
	DefaultHealthInterval = 1 * time.Minute
	DefaultConcurrency    = 10
)
