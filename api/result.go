package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r FunctionResultRequest) Valid() error {

	if r.Id == "" {
		return errors.New("request ID is required")
	}

	return nil
}

// ExecutionResult implements the REST API endpoint for retrieving the result of a function execution.
func (a *API) ExecutionResult(ctx echo.Context) error {

	// Get the request ID.
	var request FunctionResultRequest
	err := ctx.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	err = request.Valid()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("missing request ID"))
	}

	// Lookup execution result.
	result, ok := a.Node.ExecutionResult(request.Id)
	if !ok {
		return ctx.NoContent(http.StatusNotFound)
	}

	// Send the response back.
	return ctx.JSON(http.StatusOK, result)
}
