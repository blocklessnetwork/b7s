package blockless

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Executor interface {
	ExecuteFunction(requestID string, request execute.Request) (execute.Result, any, error)
}
