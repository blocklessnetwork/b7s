package api

import (
	"context"

	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
)

type Node interface {
	ExecuteFunction(ctx context.Context, req execute.Request, subgroup string) (code codes.Code, requestID string, results execute.ResultMap, peers execute.Cluster, err error)
	ExecutionResult(id string) (execute.ResultMap, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string, subgroup string) error
}
