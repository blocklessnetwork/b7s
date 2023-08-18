package mocks

import (
	"testing"

	"github.com/blocklessnetwork/b7s/models/execute"
)

type Executor struct {
	ExecFunctionFunc func(string, execute.Request) (execute.Result, error)
}

func BaselineExecutor(t *testing.T) *Executor {
	t.Helper()

	executor := Executor{
		ExecFunctionFunc: func(string, execute.Request) (execute.Result, error) {
			return GenericExecutionResult, nil
		},
	}

	return &executor
}

func (e *Executor) ExecuteFunction(requestID string, req execute.Request) (execute.Result, error) {
	return e.ExecFunctionFunc(requestID, req)
}
