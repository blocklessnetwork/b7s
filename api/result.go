package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ExecutionResult implements the REST API endpoint for retrieving the result of a function execution.
func (a *API) ExecutionResult(ctx echo.Context) error {

	// Get the request ID.
	requestID := ctx.Param("id")
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
