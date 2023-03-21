package executor

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestExecutor_RequestPaths(t *testing.T) {

	const (
		workdir        = "/var/tmp/b7s/workspace"
		requestID      = "request-id"
		functionID     = "function-id"
		functionMethod = "function-method"

		// Expected paths.
		expectedRequestWorkdir = workdir + "/t/request-id"
		expectedFSRoot         = workdir + "/t/request-id/fs"
		expectedManifestPath   = workdir + "/t/request-id/runtime-manifest.json"
		expectedEntry          = workdir + "/function-id/function-method"
	)

	executor := &Executor{
		log: mocks.NoopLogger,
		cfg: Config{
			WorkDir: workdir,
		},
	}

	// NOTE: We use filepath.Clean to have consistent separators on platforms.
	paths := executor.generateRequestPaths(requestID, functionID, functionMethod)
	assert.Equal(t, filepath.Clean(expectedRequestWorkdir), paths.workdir)
	assert.Equal(t, filepath.Clean(expectedEntry), paths.entry)
	assert.Equal(t, filepath.Clean(expectedFSRoot), paths.fsRoot)
	assert.Equal(t, filepath.Clean(expectedManifestPath), paths.manifest)
}
