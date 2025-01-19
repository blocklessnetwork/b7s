package consensus

import (
	"fmt"
	"strings"
)

// Type identifies consensus protocols suported by Bless.
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

func Parse(s string) (Type, error) {

	if s == "" {
		return 0, nil
	}

	switch strings.ToLower(s) {
	case "raft":
		return Raft, nil

	case "pbft":
		return PBFT, nil
	}

	return 0, fmt.Errorf("unknown consensus value (%s)", s)
}
