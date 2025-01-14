package main

import (
	"strings"

	"github.com/blessnetwork/b7s/models/blockless"
)

func parseNodeRole(role string) blockless.NodeRole {

	switch strings.ToLower(role) {

	case blockless.HeadNodeLabel:
		return blockless.HeadNode

	case blockless.WorkerNodeLabel:
		return blockless.WorkerNode

	default:
		panic("invalid node role specified")
	}
}
