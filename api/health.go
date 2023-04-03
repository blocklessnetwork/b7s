package api

import (
	"net/http"

	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/labstack/echo/v4"
)

// Execute implements the REST API endpoint for function execution.
func (a *API) Health(ctx echo.Context) error {

	// respond with health check	
	resp := response.Health{
		Type: "health",
		Code: http.StatusOK,
	}

	// Send the response.
	return ctx.JSON(http.StatusOK, resp)
}
