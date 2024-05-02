//go:build integration
// +build integration

package fstore_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/store/codec"
	"github.com/blocklessnetwork/b7s/testing/helpers"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

const (
	cleanupDisableEnv = "B7S_INTEG_CLEANUP_DISABLE"
)

func TestStore_InstallFunction(t *testing.T) {

	const (
		functionCID = "bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea"
		manifestURL = "https://bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea.ipfs.w3s.link/manifest.json"
		dirPattern  = "b7s-fstore-integration-test-"
	)

	// 0. Setup.

	t.Log("starting test")

	dir, err := os.MkdirTemp("", dirPattern)
	require.NoError(t, err)

	cleanupDisabled := cleanupDisabled()
	if !cleanupDisabled {
		defer os.RemoveAll(dir)
	}

	t.Logf("test dir: %v", dir)

	db := helpers.InMemoryDB(t)
	defer db.Close()

	fstore := fstore.New(mocks.NoopLogger, store.New(db, codec.NewJSONCodec()), dir)

	// 1. Function Install
	err = fstore.Install(manifestURL, functionCID)
	require.NoError(t, err)

	t.Log("function install successful")

	// 2. Verify function installation on filesystem - file structure, checksum etc.

	manifest := getManifest(t, manifestURL)

	archive := filepath.Join(dir, functionCID, manifest.Runtime.URL)
	listedChecksum, err := hex.DecodeString(manifest.Runtime.Checksum)
	require.Equal(t, listedChecksum, getChecksum(t, archive))

	t.Logf("verified checksum: checksum: %x, archive: %v", listedChecksum, archive)

	file := filepath.Join(dir, functionCID, manifest.Entry)
	info, err := os.Stat(file)
	require.NoError(t, err)
	require.NotZero(t, info.Size())

	t.Logf("verified extracted file: path: %v", file)

	// 3. Verify function record is persisted

	function, err := fstore.Get(functionCID)
	require.NoError(t, err)

	t.Log("retrieved function record")

	require.Equal(t, functionCID, function.CID)
	require.Equal(t, manifestURL, function.URL)

	// TODO: Fix manifest handling
	// require.Equal(t, manifest, function.Manifest)
	// Record has the workdir prefix trimmed.
	require.Contains(t, archive, function.Archive)
	require.Contains(t, file, function.Files)

	t.Logf("verified persisted function record")

	// 4. Verify sync functionality by deleting files and running a sync.

	require.NoError(t, os.Remove(archive))
	require.NoError(t, os.Remove(file))

	require.NoError(t, fstore.Sync(true))
	require.Equal(t, listedChecksum, getChecksum(t, archive))
	info, err = os.Stat(file)
	require.NoError(t, err)
	require.NotZero(t, info.Size())

	t.Logf("verified files reappear after sync")
}

func cleanupDisabled() bool {
	return os.Getenv(cleanupDisableEnv) == "yes"
}

func getManifest(t *testing.T, url string) blockless.FunctionManifest {
	t.Helper()

	res, err := http.Get(url)
	require.NoError(t, err)
	defer res.Body.Close()

	var manifest blockless.FunctionManifest
	err = json.NewDecoder(res.Body).Decode(&manifest)
	require.NoError(t, err)

	return manifest
}

func getChecksum(t *testing.T, path string) []byte {
	t.Helper()

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	require.NoError(t, err)

	return h.Sum(nil)
}
