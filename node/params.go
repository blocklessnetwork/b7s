package node

import (
	"errors"
	"time"
)

const (
	DefaultTopic          = "blockless/b7s/general"
	DefaultHealthInterval = 1 * time.Minute

	functionInstallTimeout = 10 * time.Second
	rollCallTimeout        = 5 * time.Second

	rollCallQueueBufferSize = 1000
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
)
