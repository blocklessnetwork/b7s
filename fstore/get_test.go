package fstore_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/fstore"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/store"
	"github.com/blocklessnetworking/b7s/testing/helpers"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestFunction_Get(t *testing.T) {

	const (
		manifestURL = "manifest.json"
		functionURL = "function.tar.gz"
		testFile    = "testdata/testFunction.tar.gz"

		testCID = "dummy-cid"
	)

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	functionPayload, err := os.ReadFile(testFile)
	require.NoError(t, err)

	hash := sha256.Sum256(functionPayload)

	// We'll create two servers, so we can link one to the other.
	msrv, fsrv := createServers(t, manifestURL, functionURL, functionPayload)

	store := store.New(helpers.InMemoryDB(t))
	fh := fstore.New(mocks.NoopLogger, store, workdir)

	address := fmt.Sprintf("%s/%v", msrv.URL, manifestURL)
	manifest, err := fh.Get(address, testCID, false)
	require.NoError(t, err)

	// Verify downloaded file.
	archive := manifest.Deployment.File
	require.FileExists(t, archive)

	ok := verifyFileHash(t, archive, hash)
	require.Truef(t, ok, "file hash does not match")

	// Shutdown both servers and retry getting the manifest - verify that the cached manifest will be returned.
	fsrv.Close()
	msrv.Close()

	_, err = fh.Get(address, testCID, true)
	require.NoError(t, err)
}

func TestFunction_GetHandlesErrors(t *testing.T) {

	const (
		manifestURL = "manifest.json"
		functionURL = "function.tar.gz"
		testFile    = "testdata/testFunction.tar.gz"

		testCID = "dummy-cid"
	)

	functionPayload, err := os.ReadFile(testFile)
	require.NoError(t, err)

	// We'll create two servers, so we can link one to the other.
	msrv, fsrv := createServers(t, manifestURL, functionURL, functionPayload)
	// NOTE: Server shutdown handled in test cases below.

	t.Run("handles failure to read manifest from store", func(t *testing.T) {

		workdir, err := os.MkdirTemp("", "b7s-function-get-")
		require.NoError(t, err)

		defer os.RemoveAll(workdir)

		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
			return mocks.GenericError
		}

		fh := fstore.New(mocks.NoopLogger, store, workdir)

		address := fmt.Sprintf("%s/%v", msrv.URL, manifestURL)
		_, err = fh.Get(address, testCID, false)
		require.Error(t, err)
	})
	t.Run("handles failure to download function", func(t *testing.T) {

		// Shutdown function server.
		fsrv.Close()

		workdir, err := os.MkdirTemp("", "b7s-function-get-")
		require.NoError(t, err)

		defer os.RemoveAll(workdir)

		store := store.New(helpers.InMemoryDB(t))
		fh := fstore.New(mocks.NoopLogger, store, workdir)

		address := fmt.Sprintf("%s/%v", msrv.URL, manifestURL)
		_, err = fh.Get(address, testCID, false)
		require.Error(t, err)
	})
	t.Run("handles failure to fetch manifest", func(t *testing.T) {

		// Shutdown manifest server.
		msrv.Close()

		workdir, err := os.MkdirTemp("", "b7s-function-get-")
		require.NoError(t, err)

		defer os.RemoveAll(workdir)

		store := store.New(helpers.InMemoryDB(t))
		fh := fstore.New(mocks.NoopLogger, store, workdir)

		address := fmt.Sprintf("%s/%v", msrv.URL, manifestURL)
		_, err = fh.Get(address, testCID, false)
		require.Error(t, err)
	})
}

func createServers(t *testing.T, manifestURL string, functionURL string, functionPayload []byte) (manifestSrv *httptest.Server, functionSrv *httptest.Server) {
	t.Helper()

	// Create function server.
	fsrv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			path := strings.TrimPrefix(req.URL.Path, "/")
			if path != functionURL {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Write(functionPayload)
		}))

	// Setup manifest that points to the function server.
	functionAddress := fmt.Sprintf("%s/%s", fsrv.URL, functionURL)
	hash := sha256.Sum256(functionPayload)
	sourceManifest := blockless.FunctionManifest{
		Deployment: blockless.Deployment{
			URI:      functionAddress,
			Checksum: fmt.Sprintf("%x", hash),
		},
	}

	// Create manifest server.
	msrv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			path := strings.TrimPrefix(req.URL.Path, "/")
			if path != manifestURL {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			payload, err := json.Marshal(sourceManifest)
			require.NoError(t, err)
			w.Write(payload)
		}))

	return msrv, fsrv
}

func verifyFileHash(t *testing.T, filename string, checksum [32]byte) bool {
	t.Helper()

	data, err := os.ReadFile(filename)
	require.NoError(t, err)

	h := sha256.Sum256(data)

	return bytes.Equal(checksum[:], h[:])
}
