package mocks

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Node implements the `Node` interface expected by the API.
type Node struct {
	ExecuteFunctionFunc        func(context.Context, execute.Request) (execute.Result, error)
	ExecutionResultFunc        func(id string) (execute.Result, bool)
	PublishFunctionInstallFunc func(ctx context.Context, uri string, cid string) error
}

func BaselineNode(t *testing.T) *Node {
	t.Helper()

	node := Node{
		ExecuteFunctionFunc: func(context.Context, execute.Request) (execute.Result, error) {
			return GenericExecutionResult, nil
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

func (n *Node) ExecuteFunction(ctx context.Context, req execute.Request) (execute.Result, error) {
	return n.ExecuteFunctionFunc(ctx, req)
}

func (n *Node) ExecutionResult(id string) (execute.Result, bool) {
	return n.ExecutionResultFunc(id)
}

func (n *Node) PublishFunctionInstall(ctx context.Context, uri string, cid string) error {
	return n.PublishFunctionInstallFunc(ctx, uri, cid)
}
