package node

import (
	"errors"
	"time"
)

const (
	DefaultTopic = "blockless/b7s/general"

	functionInstallTimeout = 10 * time.Second
	rollCallTimeout        = 5 * time.Second

	resultBufferSize = 10
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
)
