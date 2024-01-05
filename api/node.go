package api

import (
	"context"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Node interface {
	ExecuteFunction(ctx context.Context, req execute.Request, topic string) (code codes.Code, requestID string, results execute.ResultMap, peers execute.Cluster, err error)
	ExecutionResult(id string) (execute.Result, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string, topic string) error
}
