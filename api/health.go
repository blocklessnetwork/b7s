package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blessnetwork/b7s/models/response"
)

// Execute implements the REST API endpoint for function execution.
func (a *API) Health(ctx echo.Context) error {

	return ctx.JSON(
		http.StatusOK,
		response.Health{Code: http.StatusOK},
	)
}
