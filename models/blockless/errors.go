package blockless

import (
	"errors"
)

// Sentinel errors.
var (
	ErrNotFound                = errors.New("not found")
	ErrRollCallTimeout         = errors.New("roll call timed out - not enough nodes responded")
	ErrExecutionNotEnoughNodes = errors.New("not enough execution results received")
)
