package blockless

// TODO: Reconsider the package name - typically I'd use the name of the project - `b7s`.
// Package `blockless` might be too wide.

// NodeRole is a representation of the node's role.
type NodeRole uint8

// The following are all possible node roles.
const (
	HeadNode NodeRole = iota + 1
	WorkerNode
)

// The following are labels for the node roles, used when parsing the node role as a string.
const (
	HeadNodeLabel   = "head"
	WorkerNodeLabel = "worker"
)

// String returns the string representation of the node role.
func (n NodeRole) String() string {

	switch n {

	case HeadNode:
		return HeadNodeLabel
	case WorkerNode:
		return WorkerNodeLabel
	default:
		return "invalid"
	}
}
