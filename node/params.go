package node

import (
	"time"
)

const (
	DefaultTopic = "blockless/b7s/general"

	functionInstallTimeout = 10 * time.Second
	rollCallTimeout        = 5 * time.Second

	resultBufferSize = 10
)
