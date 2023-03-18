//go:build integration
// +build integration

package executor_test

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

const (
	runtimeDirEnv     = "B7S_INTEG_RUNTIME_DIR"
	cleanupDisableEnv = "B7S_INTEG_CLEANUP_DISABLE"
)

func TestExecutor_Execute(t *testing.T) {

	const (
		dirPattern = "b7s-executor-integration-test-"

		testFunction = "./testdata/md5sum/md5sum.wasm"

		functionID = "function-id"
		requestID  = "dummy-request-id"

		chunkSize = 128
		fileSize  = 256
	)

	// Setup directories.
	workspace, err := os.MkdirTemp("", dirPattern)
	require.NoError(t, err)
	if !cleanupDisabled() {
		defer os.RemoveAll(workspace)
	}

	var (
		workdir     = path.Join(workspace, "t", requestID) // request work directory
		fsRoot      = path.Join(workdir, "fs")             // function FS root
		functiondir = path.Join(workspace, functionID)     // function location
	)

	t.Logf("working directory: %v", workspace)
	createDirs(t, workdir, fsRoot, functiondir)

	// Stage executable to working directory.
	copyFunction(t, testFunction, functiondir)

	// Create a random testfile.
	testfile, hash := createTestFile(t, fsRoot, fileSize)

	// Create executor.
	executor, err := executor.New(
		mocks.NoopLogger,
		executor.WithWorkDir(workspace),
		executor.WithRuntimeDir(os.Getenv(runtimeDirEnv)),
	)
	require.NoError(t, err)

	// Execute the function.
	req := execute.Request{
		FunctionID: functionID,
		Method:     path.Base(testFunction),
		Parameters: []execute.Parameter{
			{Value: "--chunk"},
			{Value: fmt.Sprintf("%v", chunkSize)},
			{Value: "--file"},
			{Value: path.Base(testfile)}, // Specify name only because the path is relative to FS root.
		},
	}

	res, err := executor.Function(requestID, req)
	require.NoError(t, err)

	// Verify the execution result.
	require.Equal(t, response.CodeOK, res.Code)
	require.Equal(t, requestID, res.RequestID)
	require.Equal(t, hash, res.Result)

	// Verify usage info - for now, only that they are non-zero.
	cpuTimeTotal := res.Usage.CPUSysTime + res.Usage.CPUUserTime
	require.Greater(t, cpuTimeTotal, time.Duration(0))
	require.NotZero(t, res.Usage.WallClockTime)
}

func createTestFile(t *testing.T, dir string, size int) (string, string) {
	t.Helper()

	const (
		filePattern = "testfile-"
	)

	f, err := ioutil.TempFile(dir, filePattern)
	require.NoError(t, err)
	defer f.Close()

	buf := make([]byte, size)

	_, err = rand.Read(buf)
	require.NoError(t, err)

	n, err := f.Write(buf)
	require.NoError(t, err)
	require.Equal(t, size, n)

	// Calculate the hash of the file payload.
	hash := md5.New()
	_, err = hash.Write(buf)
	require.NoError(t, err)

	md5sum := fmt.Sprintf("%x", hash.Sum(nil))

	return f.Name(), md5sum
}

func createDirs(t *testing.T, dirs ...string) {
	t.Helper()

	// Create directory structure.
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		require.NoError(t, err)
	}

	return
}

func copyFunction(t *testing.T, filepath string, target string) {
	t.Helper()

	payload, err := os.ReadFile(filepath)
	require.NoError(t, err)

	_, name := path.Split(filepath)
	targetPath := path.Join(target, name)

	err = os.WriteFile(targetPath, payload, os.ModePerm)
	require.NoError(t, err)
}

func cleanupDisabled() bool {
	return os.Getenv(cleanupDisableEnv) == "yes"
}
