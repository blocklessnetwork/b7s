package node

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/models/api"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
)

// TODO: Unify style on all JSON models - use snake_case.

func (n *Node) HandleRequestExecute(ctx echo.Context) error {

	// TODO: Uses different models from before. Do we need specific model for
	// REST for this?
	var req request.Execute
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	execReq := execute.Request{
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
	}

	response, err := n.execute.Function(execReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("execution failed: %w", err))
	}

	// Determine which status code to send.
	code := http.StatusOK
	if err != nil {
		code = http.StatusInternalServerError
	}

	return ctx.JSON(code, response)
}

func (n *Node) HandleFunctionInstall(ctx echo.Context) error {
	return fmt.Errorf("TBD: Not implemented")
}

func (n *Node) HandleGetExecuteResponse(ctx echo.Context) error {

	// Unpack the request.
	var req api.RequestExecuteResponse
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	// Lookup this execution response.
	res, ok := n.excache.Get(req.ID)
	if !ok {
		return ctx.NoContent(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, res)
}
