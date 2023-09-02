package consensus

import (
	"fmt"
)

// Type identifies consensus protocols suported by Blockless.
type Type uint

const (
	Raft Type = iota + 1
	PBFT
)

func (t Type) String() string {
	switch t {
	case Raft:
		return "Raft"
	case PBFT:
		return "PBFT"
	default:
		return fmt.Sprintf("unknown: %d", t)
	}
}

func (t Type) Valid() bool {
	switch t {
	case Raft, PBFT:
		return true
	default:
		return false
	}
}
