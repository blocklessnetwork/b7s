package api

import (
	"context"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

type Node interface {
	ExecuteFunction(context.Context, execute.Request) (code codes.Code, requestID string, results execute.ResultMap, peers execute.Cluster, err error)
	ExecutionResult(id string) (execute.Result, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string) error
}
