package main

import (
	"strings"

	"github.com/blessnetwork/b7s/models/bls"
)

func parseNodeRole(role string) bls.NodeRole {

	switch strings.ToLower(role) {

	case bls.HeadNodeLabel:
		return bls.HeadNode

	case bls.WorkerNodeLabel:
		return bls.WorkerNode

	default:
		panic("invalid node role specified")
	}
}
