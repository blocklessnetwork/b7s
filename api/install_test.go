package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/api"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestAPI_FunctionInstall(t *testing.T) {
	t.Run("nominal case", func(t *testing.T) {
		t.Parallel()

		req := api.InstallFunctionRequest{
			URI: "dummy-function-id",
			CID: "dummy-cid",
		}

		srv := setupAPI(t)

		rec, ctx, err := setupRecorder(installEndpoint, req)
		require.NoError(t, err)

		err = srv.Install(ctx)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	})
}

func TestAPI_FunctionInstall_HandlesErrors(t *testing.T) {
	t.Run("missing URI and CID", func(t *testing.T) {
		t.Parallel()

		req := api.InstallFunctionRequest{
			URI: "",
			CID: "",
		}

		srv := setupAPI(t)

		_, ctx, err := setupRecorder(installEndpoint, req)
		require.NoError(t, err)

		err = srv.Install(ctx)
		require.Error(t, err)

		echoErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)

		require.Equal(t, http.StatusBadRequest, echoErr.Code)
	})
	t.Run("node install takes too long", func(t *testing.T) {
		t.Parallel()

		const (
			// The API times out after 10 seconds.
			installDuration = 11 * time.Second
		)

		node := mocks.BaselineNode(t)
		node.PublishFunctionInstallFunc = func(context.Context, string, string) error {
			time.Sleep(installDuration)
			return nil
		}

		req := api.InstallFunctionRequest{
			URI: "dummy-uri",
			CID: "dummy-cid",
		}

		srv := api.New(mocks.NoopLogger, node)

		rec, ctx, err := setupRecorder(installEndpoint, req)
		require.NoError(t, err)

		err = srv.Install(ctx)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, rec.Result().StatusCode)

		var res api.InstallFunctionResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

		num, err := strconv.Atoi(res.Code)
		require.NoError(t, err)

		require.Equal(t, http.StatusRequestTimeout, num)
	})
	t.Run("node fails to install function", func(t *testing.T) {
		t.Parallel()

		node := mocks.BaselineNode(t)
		node.PublishFunctionInstallFunc = func(context.Context, string, string) error {
			return mocks.GenericError
		}

		srv := api.New(mocks.NoopLogger, node)

		req := api.InstallFunctionRequest{
			URI: "dummy-uri",
			CID: "dummy-cid",
		}

		_, ctx, err := setupRecorder(installEndpoint, req)
		require.NoError(t, err)

		err = srv.Install(ctx)
		require.Error(t, err)

		echoErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)

		require.Equal(t, http.StatusInternalServerError, echoErr.Code)
	})
}

func TestAPI_InstallFunction_HandlesMalformedRequests(t *testing.T) {

	srv := setupAPI(t)

	const (
		wrongFieldType = `
		{
			"uri": "dummy-uri",
			"cid": 14
		}`

		unclosedBracket = `
		{
			"uri": "dummy-uri",
			"cid": "dummy-cid"
		`

		validJSON = `
		{
			"uri": "dummy-uri",
			"cid": "dummy-cid"
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

			_, ctx, err := setupRecorder(installEndpoint, test.payload, prepare)
			require.NoError(t, err)

			err = srv.Install(ctx)
			require.Error(t, err)

			echoErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)

			require.Equal(t, http.StatusBadRequest, echoErr.Code)
		})
	}
}
