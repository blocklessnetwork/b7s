package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// API provides REST API functionality for the Blockless head node.
type API struct {
	log  zerolog.Logger
	srv  *echo.Echo
	node Node
}

// New creates a new instance of a Blockless head node REST API. Access to node data is provided by the provided `node`.
func New(log zerolog.Logger, srv *echo.Echo, node Node) *API {

	api := API{
		log:  log.With().Str("component", "api").Logger(),
		srv:  srv,
		node: node,
	}

	api.registerRoutes()

	return &api
}

// registerRoutes sets the endpoint handlers.
func (a *API) registerRoutes() {
	a.srv.GET("/api/v1/health", a.Health)
	a.srv.POST("/api/v1/functions/execute", a.Execute)
	a.srv.POST("/api/v1/functions/install", a.Install)
	a.srv.POST("/api/v1/functions/requests/result", a.ExecutionResult)
}
