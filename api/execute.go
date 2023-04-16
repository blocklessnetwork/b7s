package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/models/api/request"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute implements the REST API endpoint for function execution.
func (a *API) Execute(ctx echo.Context) error {

	// Unpack the API request.
	var req request.Execute
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	// TODO: Check - We perhaps want to return the request ID and not wait for the execution, right?
	// It's probable that it will time out anyway, right?

	// Get the execution results.
	results, err := a.node.ExecuteFunction(ctx.Request().Context(), execute.Request(req))
	_ = results

	// TODO: Correct the API response.
	// Create the API response.
	// res := response.Execute{
	// 	Code:      results.Code,
	// 	RequestID: results.RequestID,
	// 	Result:    results.Result.Stdout,
	// 	ResultEx:  results.Result,
	// }

	// Send the response.
	return ctx.JSON(http.StatusOK, results)
}
