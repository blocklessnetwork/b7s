package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/api"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestAPI_Execute(t *testing.T) {

	executionResult := execute.Result{
		Result: execute.RuntimeOutput{
			Stdout:   "dummy-failed-execution-result",
			Stderr:   "dummy-failed-execution-log",
			ExitCode: 0,
		},
	}
	peerIDs := []peer.ID{
		mocks.GenericPeerID,
	}
	expectedCode := codes.OK

	node := mocks.BaselineNode(t)
	node.ExecuteFunctionFunc = func(context.Context, execute.Request, string) (codes.Code, string, execute.ResultMap, execute.Cluster, error) {

		res := execute.ResultMap{
			mocks.GenericPeerID: executionResult,
		}

		cluster := execute.Cluster{
			Peers: peerIDs,
		}

		return expectedCode, mocks.GenericUUID.String(), res, cluster, nil
	}

	srv := api.New(mocks.NoopLogger, node)

	req := mocks.GenericExecutionRequest

	rec, ctx, err := setupRecorder(executeEndpoint, req)
	require.NoError(t, err)

	err = srv.ExecuteFunction(ctx)
	require.NoError(t, err)

	var res api.ExecutionResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	require.Equal(t, expectedCode.String(), res.Code)

	require.Len(t, res.Cluster.Peers, 1)
	require.Equal(t, res.Cluster.Peers, peerIDs)

	require.Len(t, res.Results, 1)
	require.Equal(t, executionResult.Result, res.Results[0].Result)
	require.Equal(t, float64(100), res.Results[0].Frequency)
	require.Equal(t, peerIDs, res.Results[0].Peers)

	require.Equal(t, mocks.GenericUUID.String(), res.RequestId)
}

func TestAPI_Execute_HandlesErrors(t *testing.T) {

	executionResult := execute.Result{
		Result: execute.RuntimeOutput{
			Stdout:   "dummy-failed-execution-result",
			Stderr:   "dummy-failed-execution-log",
			ExitCode: 1,
		},
	}

	expectedCode := codes.Error

	node := mocks.BaselineNode(t)
	node.ExecuteFunctionFunc = func(context.Context, execute.Request, string) (codes.Code, string, execute.ResultMap, execute.Cluster, error) {

		res := execute.ResultMap{
			mocks.GenericPeerID: executionResult,
		}

		return expectedCode, "", res, execute.Cluster{}, mocks.GenericError
	}

	srv := api.New(mocks.NoopLogger, node)

	req := mocks.GenericExecutionRequest

	rec, ctx, err := setupRecorder(executeEndpoint, req)
	require.NoError(t, err)

	err = srv.ExecuteFunction(ctx)
	require.NoError(t, err)

	var res api.ExecutionResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, expectedCode.String(), res.Code)

	require.Len(t, res.Results, 1)
	require.Equal(t, executionResult.Result, res.Results[0].Result)
	require.Equal(t, float64(100), res.Results[0].Frequency)
	require.Len(t, res.Results[0].Peers, 1)
	require.Equal(t, mocks.GenericPeerID, res.Results[0].Peers[0])
}

func TestAPI_Execute_HandlesMalformedRequests(t *testing.T) {

	api := setupAPI(t)

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

			err = api.ExecuteFunction(ctx)
			require.Error(t, err)

			echoErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)

			require.Equal(t, http.StatusBadRequest, echoErr.Code)
		})
	}
}
