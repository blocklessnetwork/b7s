package blockless

import (
	"context"

	"github.com/blessnetwork/b7s/models/execute"
)

type Executor interface {
	ExecuteFunction(ctx context.Context, requestID string, request execute.Request) (execute.Result, error)
}
