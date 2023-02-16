package executor

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
)

// Function will execute the Blockless function defined by the execution request.
func (e *Executor) Function(req execute.Request) (execute.Result, error) {

	// Get a new request ID.
	uuid, err := uuid.NewRandom()
	if err != nil {
		// Should NEVER really happen.
		res := execute.Result{
			Code:      response.CodeError,
			RequestID: "",
			Result:    "Could not generate request ID",
		}

		return res, fmt.Errorf("could not generate request ID: %w", err)
	}
	requestID := uuid.String()

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
