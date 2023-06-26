package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// ExecuteRequest describes the payload for the REST API request for function execution.
type ExecuteRequest execute.Request

// ExecuteResponse describes the REST API response for function execution.
type ExecuteResponse struct {
	Code      codes.Code    `json:"code,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
	Message   string        `json:"message,omitempty"`
	Result    ExecuteResult `json:"results,omitempty"` // TODO: (raft) On main - map[string]ExecuteResult
}

// ExecuteResult represents the API representation of a single execution response.
// It is similar to the model in `execute.Result`, except it omits the usage information for now.
type ExecuteResult struct {
	Code      codes.Code            `json:"code,omitempty"`
	Result    execute.RuntimeOutput `json:"result,omitempty"`
	RequestID string                `json:"request_id,omitempty"`
}

// Execute implements the REST API endpoint for function execution.
func (a *API) Execute(ctx echo.Context) error {

	// Unpack the API request.
	var req ExecuteRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	// TODO: Check - We perhaps want to return the request ID and not wait for the execution, right?
	// It's probable that it will time out anyway, right?

	// Get the execution result.
	code, result, err := a.node.ExecuteFunction(ctx.Request().Context(), execute.Request(req))
	if err != nil {
		a.log.Warn().Str("function_id", req.FunctionID).Err(err).Msg("node failed to execute function")
	}

	// Transform the node response format to the one returned by the API.
	res := ExecuteResponse{
		Code:      code,
		RequestID: result.RequestID,
		Result: ExecuteResult{
			Code:      result.Code,
			Result:    result.Result,
			RequestID: result.RequestID,
		},
	}

	// Communicate the reason for failure in these cases.
	if errors.Is(err, blockless.ErrRollCallTimeout) || errors.Is(err, blockless.ErrExecutionNotEnoughNodes) {
		res.Message = err.Error()
	}

	// Send the response.
	return ctx.JSON(http.StatusOK, res)
}
