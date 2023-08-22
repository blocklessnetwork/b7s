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
		return "pBFT"
	default:
		return fmt.Sprintf("unknown: %d", t)
	}
}
