package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/blocklessnetworking/b7s/api"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

const (
	executeEndpoint = "/api/v1/functions/execute"
	installEndpoint = "/api/v1/functions/install"
	resultEndpoint  = "/api/v1/functions/requests/result"
)

func setupAPI(t *testing.T) *api.API {
	t.Helper()

	var (
		logger = mocks.NoopLogger
		node   = mocks.BaselineNode(t)
	)

	api := api.New(logger, node)

	return api
}

func setupRecorder(endpoint string, input interface{}, options ...func(*http.Request)) (*httptest.ResponseRecorder, echo.Context, error) {

	payload, ok := input.([]byte)
	if !ok {
		var err error
		payload, err = json.Marshal(input)
		if err != nil {
			return nil, echo.New().AcquireContext(), fmt.Errorf("could not encode input: %w", err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	for _, opt := range options {
		opt(req)
	}

	rec := httptest.NewRecorder()

	ctx := echo.New().NewContext(req, rec)

	return rec, ctx, nil
}
