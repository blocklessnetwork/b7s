package mocks

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Node implements the `Node` interface expected by the API.
type Node struct {
	ExecuteFunctionFunc        func(context.Context, execute.Request) (codes.Code, map[string]execute.Result, error)
	ExecutionResultFunc        func(id string) (execute.Result, bool)
	PublishFunctionInstallFunc func(ctx context.Context, uri string, cid string) error
}

func BaselineNode(t *testing.T) *Node {
	t.Helper()

	node := Node{
		ExecuteFunctionFunc: func(context.Context, execute.Request) (codes.Code, map[string]execute.Result, error) {

			res := map[string]execute.Result{
				GenericPeerID.String(): GenericExecutionResult,
			}

			return GenericExecutionResult.Code, res, nil
		},
		ExecutionResultFunc: func(id string) (execute.Result, bool) {
			return GenericExecutionResult, true
		},
		PublishFunctionInstallFunc: func(ctx context.Context, uri string, cid string) error {
			return nil
		},
	}

	return &node
}

func (n *Node) ExecuteFunction(ctx context.Context, req execute.Request) (codes.Code, map[string]execute.Result, error) {
	return n.ExecuteFunctionFunc(ctx, req)
}

func (n *Node) ExecutionResult(id string) (execute.Result, bool) {
	return n.ExecutionResultFunc(id)
}

func (n *Node) PublishFunctionInstall(ctx context.Context, uri string, cid string) error {
	return n.PublishFunctionInstallFunc(ctx, uri, cid)
}
