package node

import (
	"errors"
	"time"
)

const (
	DefaultTopic           = "blockless/b7s/general"
	DefaultHealthInterval  = 1 * time.Minute
	DefaultRollCallTimeout = 5 * time.Second
	DefaultConcurrency     = 10

	functionInstallTimeout = 10 * time.Second

	rollCallQueueBufferSize = 1000

	syncInterval = time.Hour
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	errRollCallTimeout    = errors.New("roll call timed out")
)
