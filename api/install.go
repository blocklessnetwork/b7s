package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (r FunctionInstallRequest) Valid() error {

	if r.Cid == "" {
		return errors.New("function CID is required")
	}

	return nil
}

func (a *API) InstallFunction(ctx echo.Context) error {

	// Unpack the API request.
	var req FunctionInstallRequest
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	err = req.Valid()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
	}

	err = a.Node.PublishFunctionInstall(ctx.Request().Context(), req.Uri, req.Cid, req.Topic)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("function installation failed: %w", err))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": strconv.Itoa(http.StatusOK),
	})

}
