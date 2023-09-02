package raft

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

// TODO: Have this as a global definition somewhere.
type Executor interface {
	ExecuteFunction(requestID string, req execute.Request) (execute.Result, error)
}
