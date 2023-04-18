package node

import (
	"errors"
	"time"
)

const (
	DefaultTopic            = "blockless/b7s/general"
	DefaultHealthInterval   = 1 * time.Minute
	DefaultRollCallTimeout  = 5 * time.Second
	DefaultExecutionTimeout = 10 * time.Second
	DefaultConcurrency      = 10

	rollCallQueueBufferSize = 1000

	syncInterval = time.Hour
)

var (
	ErrUnsupportedMessage      = errors.New("unsupported message")
	errRollCallTimeout         = errors.New("roll call timed out - not enough nodes responded")
	errExecutionNotEnoughNodes = errors.New("not enough execution results received")
)
