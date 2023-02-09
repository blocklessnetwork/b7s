package executor

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
)

// execute will execute the Blockless function defined by the execution request.
func (e *Executor) Function(req execute.Request) (execute.Response, error) {

	// Get a new execution ID.
	uuid, err := uuid.NewRandom()
	if err != nil {
		// Should NEVER really happen.
		res := execute.Response{
			Code:      response.CodeError,
			RequestID: "",
			Result:    "Could not generate execution ID",
		}

		return res, fmt.Errorf("could not generate execution ID: %w", err)
	}
	executionID := uuid.String()

	// Execute the function.
	out, err := e.execute(executionID, req)
	if err != nil {

		res := execute.Response{
			Code:      response.CodeError,
			RequestID: executionID,
		}

		return res, fmt.Errorf("function execution failed: %w", err)
	}

	// TODO: Execution response memory store.

	res := execute.Response{
		Code:      response.CodeOK,
		RequestID: executionID,
		Result:    out,
	}

	return res, nil
}
