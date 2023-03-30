package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/models/api/request"
)

const (
	functionInstallTimeout = 10 * time.Second
)

func (a *API) Install(ctx echo.Context) error {

	// Unpack the API request.
	var req request.InstallFunction
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("could not unpack request: %w", err))
	}

	if req.URI == "" && req.CID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("URI or CID are required"))
	}

	// Add a deadline to the context.
	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), functionInstallTimeout)
	defer cancel()

	// Start function install in a separate goroutine and signal when it's done.
	fnErr := make(chan error)
	go func() {
		err = a.node.FunctionInstall(reqCtx, req.URI, req.CID)
		fnErr <- err
	}()

	// Wait until either function install finishes, or request times out.
	select {

	// Context timed out.
	case <-reqCtx.Done():

		status := http.StatusRequestTimeout
		if !errors.Is(reqCtx.Err(), context.DeadlineExceeded) {
			status = http.StatusInternalServerError
		}

		// return inner code as body
		return ctx.JSON(200, map[string]interface{}{
			"code": strconv.Itoa(status),
		})

	// Work done.
	case err = <-fnErr:
		break
	}

	// Check if function install succeeded and handle error or return response.
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("function installation failed: %w", err))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": strconv.Itoa(http.StatusOK),
	})

}
