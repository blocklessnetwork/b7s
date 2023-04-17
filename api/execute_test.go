package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/api"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestAPI_Execute(t *testing.T) {

	srv := setupAPI(t)

	req := mocks.GenericExecutionRequest

	rec, ctx, err := setupRecorder(executeEndpoint, req)
	require.NoError(t, err)

	err = srv.Execute(ctx)
	require.NoError(t, err)

	var res api.ExecuteResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, mocks.GenericExecutionResult.Code, res.Code)
	require.Equal(t, mocks.GenericExecutionResult.RequestID, res.RequestID)
	require.Equal(t, mocks.GenericExecutionResult.Result.Stdout, res.Result)
	require.Equal(t, mocks.GenericExecutionResult.Result, res.ResultEx)
}

func TestAPI_Execute_HandlesErrors(t *testing.T) {

	executionResult := execute.Result{
		Result: execute.RuntimeOutput{
			Stdout:   "dummy-failed-execution-result",
			Stderr:   "dummy-failed-execution-log",
			ExitCode: 1,
		},
	}

	node := mocks.BaselineNode(t)
	node.ExecuteFunctionFunc = func(context.Context, execute.Request) (execute.Result, error) {
		return executionResult, mocks.GenericError
	}

	srv := api.New(mocks.NoopLogger, node)

	req := mocks.GenericExecutionRequest

	rec, ctx, err := setupRecorder(executeEndpoint, req)
	require.NoError(t, err)

	err = srv.Execute(ctx)
	require.NoError(t, err)

	var res api.ExecuteResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	require.Equal(t, executionResult.Code, res.Code)
	require.Equal(t, executionResult.Result.Stdout, res.Result)
	require.Equal(t, executionResult.Result, res.ResultEx)
}

func TestAPI_Execute_HandlesMalformedRequests(t *testing.T) {

	api := setupAPI(t)
	_ = api

	const (
		wrongFieldType = `
		{
			"function_id" : "generic-function-id",
			"method" : 14,
			"parameters" : [ {"name":"generic-param-name","value":"generic-param-value"} ]
		}`

		unclosedBracket = `
		{
			"function_id" : "generic-function-id",
			"method" : "wasm",
			"parameters" : [ {"name":"generic-param-name","value":"generic-param-value"} ]
		`

		validJSON = `
		{
			"function_id" : "generic-function-id",
			"method" : "wasm",
			"parameters" : [ {"name":"generic-param-name","value":"generic-param-value"} ]
		}`
	)

	tests := []struct {
		name        string
		payload     []byte
		contentType string
	}{
		{
			name:        "wrong field type",
			payload:     []byte(wrongFieldType),
			contentType: echo.MIMEApplicationJSON,
		},
		{
			name:        "malformed JSON",
			payload:     []byte(unclosedBracket),
			contentType: echo.MIMEApplicationJSON,
		},
		{
			name:    "valid JSON with no MIME type",
			payload: []byte(validJSON),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			prepare := func(req *http.Request) {
				req.Header.Set(echo.HeaderContentType, test.contentType)
			}

			_, ctx, err := setupRecorder(executeEndpoint, test.payload, prepare)
			require.NoError(t, err)

			err = api.Execute(ctx)
			require.Error(t, err)

			echoErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)

			require.Equal(t, http.StatusBadRequest, echoErr.Code)
		})
	}
}
