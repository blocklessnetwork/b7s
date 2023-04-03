package node

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

type Executor interface {
	ExecuteFunction(requestID string, req execute.Request) (execute.Result, error)
}
