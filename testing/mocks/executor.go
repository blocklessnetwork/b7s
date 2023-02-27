package mocks

import (
	"testing"

	"github.com/blocklessnetworking/b7s/models/execute"
)

type Executor struct {
	ExecFunctionFunc func(execute.Request) (execute.Result, error)
}

func BaselineExecutor(t *testing.T) *Executor {
	t.Helper()

	executor := Executor{
		ExecFunctionFunc: func(execute.Request) (execute.Result, error) {
			return GenericExecutionResult, nil
		},
	}

	return &executor
}

func (e *Executor) Function(req execute.Request) (execute.Result, error) {
	return e.ExecFunctionFunc(req)
}
