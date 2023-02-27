package function

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/store"
	"github.com/blocklessnetworking/b7s/testing/helpers"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestFunction_GetJSON(t *testing.T) {

	var (
		workdir  = "/"
		manifest = blockless.FunctionManifest{
			ID:          "generic-id",
			Name:        "generic-name",
			Description: "generic-description",
			Function: blockless.Function{
				ID:      "function-id",
				Name:    "function-name",
				Runtime: "generic-runtime",
			},
			Deployment: blockless.Deployment{
				CID:      "generic-cid",
				Checksum: "1234567890",
				URI:      "generic-uri",
			},
			FSRootPath: "/var/tmp/blockless/",
			Entry:      "/var/tmp/blockless/app.wasm",
		}
	)

	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			payload, err := json.Marshal(manifest)
			require.NoError(t, err)
			w.Write(payload)
		}))
	defer srv.Close()

	store := store.New(helpers.InMemoryDB(t))
	fh := NewHandler(mocks.NoopLogger, store, workdir)

	var downloaded blockless.FunctionManifest
	err := fh.getJSON(srv.URL, &downloaded)
	require.NoError(t, err)

	require.Equal(t, manifest, downloaded)
}

func TestFunction_GetJSONHandlesErrors(t *testing.T) {

	const (
		workdir = "/"
	)

	tests := []struct {
		name string

		statusCode int
		payload    []byte
	}{
		{
			name: "handles malformed JSON",

			statusCode: http.StatusOK,
			// JSON payload without closing brace.
			payload: []byte(`{
			"id":"generic-id",
			"name":"generic-name",
			"description":"generic-description",
			"function": {
				"id":"function-id",
				"name":"function-name",
				"runtime":"generic-runtime"
			},
			"deployment":{
				"cid":"generic-cid",
				"checksum":"1234567890",
				"uri":"generic-uri"
			},
			"runtime":{},
			"fs_root_path":"/var/tmp/blockless/",
			"entry":"/var/tmp/blockless/app.wasm"`), // <- missing closing brace
		},
		{
			name: "handles unexpected format",

			statusCode: http.StatusOK,
			// Valid JSON payload but wrong format - number instead of a textual fiel.d
			payload: []byte(`{
				"id":"generic-id",
				"name":"generic-name",
				"description":"generic-description",
				"function": {
					"id":"function-id",
					"name":"function-name",
					"runtime":"generic-runtime"
				},
				"deployment":{
					"cid":"generic-cid",
					"checksum":"1234567890",
					"uri":"generic-uri"
				},
				"runtime":{},
				"fs_root_path":"/var/tmp/blockless/",
				"entry":999
			}`),
		},
		{
			name:       "handles missing data",
			statusCode: http.StatusInternalServerError,
			payload:    []byte{},
		},
	}

	for _, test := range tests {

		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(test.statusCode)
					w.Write(test.payload)
				}))
			defer srv.Close()

			store := store.New(helpers.InMemoryDB(t))
			fh := NewHandler(mocks.NoopLogger, store, workdir)

			var response blockless.FunctionManifest
			err := fh.getJSON(srv.URL, &response)
			require.Error(t, err)
		})
	}
}
