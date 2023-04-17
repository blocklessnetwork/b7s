package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ExecutionResultRequest describes the payload for the REST API request for execution result.
type ExecutionResultRequest struct {
	ID string `json:"id"`
}

// ExecutionResult implements the REST API endpoint for retrieving the result of a function execution.
func (a *API) ExecutionResult(ctx echo.Context) error {

	// Get the request ID.
	var request ExecutionResultRequest
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	requestID := request.ID
	if requestID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("missing request ID"))
	}

	// Lookup execution result.
	result, ok := a.node.ExecutionResult(requestID)
	if !ok {
		return ctx.NoContent(http.StatusNotFound)
	}

	// Send the response back.
	return ctx.JSON(http.StatusOK, result)
}
