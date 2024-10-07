package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ExecutionResult implements the REST API endpoint for retrieving the result of a function execution.
func (a *API) ExecutionResult(ctx echo.Context) error {

	// Get the request ID.
	var request FunctionResultRequest
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	requestID := request.Id
	if requestID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("missing request ID"))
	}

	// Lookup execution result.
	result, ok := a.Node.ExecutionResult(requestID)
	if !ok {
		return ctx.NoContent(http.StatusNotFound)
	}

	// TODO: Output format fixed.

	// Send the response back.
	return ctx.JSON(http.StatusOK, result)
}
