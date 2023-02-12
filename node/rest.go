package node

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/models/api"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
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

	// TODO: Broken - if we have REST, we're a head node and we're not executing stuff directly.
	// Fix this.

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

	// Unpack the request.
	var req api.RequestFunctionInstall
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
		err = n.functionInstall(reqCtx, req.URI, req.CID)
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
		return ctx.NoContent(status)

	// Work done.
	case err = <-fnErr:
		break
	}

	// Check if function install succeeded and handle error or return response.
	if err != nil {
		n.log.Error().
			Err(err).
			Str("uri", req.URI).
			Str("cid", req.CID).
			Msg("failed to install function")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	// Create response.
	res := api.ResponseFunctionInstall{
		Code: response.CodeOK,
	}

	return ctx.JSON(http.StatusOK, res)
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

// TODO: Move these functions.

func (n *Node) functionInstall(ctx context.Context, uri string, cid string) error {

	var req request.InstallFunction
	if uri != "" {
		var err error
		req, err = createInstallMessageFromURI(uri)
		if err != nil {
			return fmt.Errorf("could not create install message from URI: %W", err)
		}
	} else {
		req = createInstallMessageFromCID(cid)
	}

	n.log.Info().
		Str("url", req.ManifestURL).
		Msg("Requesting to message peer for function installation")

	return nil
}

// createInstallMessageFromURI creates a MsgInstallFunction from the given URI.
// CID is calculated as a SHA-256 hash of the URI.
func createInstallMessageFromURI(uri string) (request.InstallFunction, error) {

	cid, err := deriveCIDFromURI(uri)
	if err != nil {
		return request.InstallFunction{}, fmt.Errorf("could not determine cid: %w", err)
	}

	msg := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		ManifestURL: uri,
		CID:         cid,
	}

	return msg, nil
}

// createInstallMessageFromCID creates the MsgInstallFunction from the given CID.
func createInstallMessageFromCID(cid string) request.InstallFunction {

	req := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		ManifestURL: fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid),
		CID:         cid,
	}

	return req
}

func deriveCIDFromURI(uri string) (string, error) {

	h := sha256.New()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return "", fmt.Errorf("could not calculate hash: %w", err)
	}
	cid := fmt.Sprintf("%x", h.Sum(nil))

	return cid, nil
}
