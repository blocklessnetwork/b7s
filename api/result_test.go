package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/api"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestAPI_ExecutionResult(t *testing.T) {
	t.Run("nominal case", func(t *testing.T) {
		t.Parallel()

		srv := setupAPI(t)

		req := api.FunctionResultRequest{
			Id: mocks.GenericString,
		}

		rec, ctx, err := setupRecorder(resultEndpoint, req)
		require.NoError(t, err)

		err = srv.ExecutionResult(ctx)
		require.NoError(t, err)

		var res execute.ResultMap
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
		require.Equal(t, mocks.GenericExecutionResultMap, res)
	})
	t.Run("response not found", func(t *testing.T) {

		node := mocks.BaselineNode(t)
		node.ExecutionResultFunc = func(id string) (execute.ResultMap, bool) {
			return execute.ResultMap{}, false
		}

		srv := api.New(mocks.NoopLogger, node)

		req := api.FunctionResultRequest{
			Id: "dummy-request-id",
		}

		rec, ctx, err := setupRecorder(resultEndpoint, req)
		require.NoError(t, err)

		err = srv.ExecutionResult(ctx)
		require.NoError(t, err)

		require.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	})
}

func TestAPI_ExecutionResult_HandlesErrors(t *testing.T) {

	srv := setupAPI(t)

	const (
		emptyIDPayload = `
		{
			"id": ""
		}`

		wrongFieldType = `
		{
			"id": 14
		}`

		unclosedBracket = `
		{
			"id": "dummy-id",
		`

		validJSON = `
		{
			"id": "dummy-id"
		}`
	)

	tests := []struct {
		name           string
		payload        []byte
		contentType    string
		expectedStatus int
	}{
		{
			name:           "empty request ID",
			payload:        []byte(emptyIDPayload),
			contentType:    echo.MIMEApplicationJSON,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong field type",
			payload:        []byte(wrongFieldType),
			contentType:    echo.MIMEApplicationJSON,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "malformed JSON",
			payload:        []byte(unclosedBracket),
			contentType:    echo.MIMEApplicationJSON,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid JSON with no MIME type set",
			payload:        []byte(validJSON),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			prepare := func(req *http.Request) {
				req.Header.Set(echo.HeaderContentType, test.contentType)
			}

			_, ctx, err := setupRecorder(resultEndpoint, test.payload, prepare)
			require.NoError(t, err)

			err = srv.ExecutionResult(ctx)
			require.Error(t, err)

			echoErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)

			require.Equal(t, test.expectedStatus, echoErr.Code)
		})
	}
}
