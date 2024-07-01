package api

import (
	"context"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
)

type Node interface {
	ExecuteFunction(ctx context.Context, req execute.Request, subgroup string) (code codes.Code, requestID string, results response.ExecutionResultMap, peers execute.Cluster, err error)
	ExecutionResult(id string) (execute.Result, bool)
	PublishFunctionInstall(ctx context.Context, uri string, cid string, subgroup string) error
}
