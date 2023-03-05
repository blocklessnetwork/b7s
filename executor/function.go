package executor

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
)

// Function will execute the Blockless function defined by the execution request.
func (e *Executor) Function(requestID string, req execute.Request) (execute.Result, error) {

	// Execute the function.
	out, err := e.execute(requestID, req)
	if err != nil {

		res := execute.Result{
			Code:      response.CodeError,
			RequestID: requestID,
		}

		return res, fmt.Errorf("function execution failed: %w", err)
	}

	res := execute.Result{
		Code:      response.CodeOK,
		RequestID: requestID,
		Result:    out,
	}

	return res, nil
}
