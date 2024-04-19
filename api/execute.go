package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/node/aggregate"
)

// Execute implements the REST API endpoint for function execution.
func (a *API) Execute(ctx echo.Context) error {

	// Unpack the API request.
	var req ExecuteRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	// Get the execution result.
	code, id, results, cluster, err := a.Node.ExecuteFunction(ctx.Request().Context(), execute.Request(req.Request), req.Topic)
	if err != nil {
		a.Log.Warn().Str("function", req.FunctionID).Err(err).Msg("node failed to execute function")
	}

	// Transform the node response format to the one returned by the API.
	res := ExecuteResponse{
		Code:      code,
		RequestID: id,
		Results:   aggregate.Aggregate(results),
		Cluster:   cluster,
	}

	// Communicate the reason for failure in these cases.
	if errors.Is(err, blockless.ErrRollCallTimeout) || errors.Is(err, blockless.ErrExecutionNotEnoughNodes) {
		res.Message = err.Error()
	}

	// Send the response.
	return ctx.JSON(http.StatusOK, res)
}
