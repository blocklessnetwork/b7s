package api

import (
	"context"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

type Node interface {
	ExecuteFunction(context.Context, execute.Request) (codes.Code, map[string]execute.Result, error)
	ExecutionResult(id string) (execute.Result, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string) error
}
