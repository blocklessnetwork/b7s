package main

import (
	"errors"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func parseNodeRole(role string) (blockless.NodeRole, error) {

	switch role {

	case blockless.HeadNodeLabel:
		return blockless.HeadNode, nil

	case blockless.WorkerNodeLabel:
		return blockless.WorkerNode, nil

	default:
		return 0, errors.New("invalid node role")
	}
}
