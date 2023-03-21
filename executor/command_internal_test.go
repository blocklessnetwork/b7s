package executor

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestExecute_CreateCMD(t *testing.T) {

	var (
		runtimeDir     = "/usr/local/bin"
		workdir        = "/var/tmp/b7s"
		functionID     = "function-id"
		functionMethod = "function-method"

<<<<<<< HEAD
		executablePath = filepath.Join(runtimeDir, blockless.RuntimeCLI())
=======
		executablePath = filepath.Join(runtimeDir, blockless.RuntimeCLI)
>>>>>>> 6110937 (Runtime CLI name fixed on Windows)

		requestID   = mocks.GenericUUID.String()
		stdin       = "dummy stdin payload"
		environment = getEnvVars(t)

		request = execute.Request{
			Config: execute.Config{
				Stdin:       &stdin,
				Environment: environment,
			},
		}
	)

	executor := Executor{
		log: mocks.NoopLogger,
		cfg: Config{
			RuntimeDir:     runtimeDir,
			WorkDir:        workdir,
			ExecutableName: blockless.RuntimeCLI(),
		},
	}
	paths := executor.generateRequestPaths(requestID, functionID, functionMethod)

	// Create command.
	cmd := executor.createCmd(paths, request)
	require.NotNil(t, cmd)

	// Verify command to be executed is correct.
	require.Equal(t, executablePath, cmd.Path)

	// Verify CLI arguments are correct.
	require.Len(t, cmd.Args, 2)
	require.Equal(t, executablePath, cmd.Args[0])
	require.Equal(t, paths.manifest, cmd.Args[1])

	// Verify working directory is correct.
	require.Equal(t, paths.workdir, cmd.Dir)

	// Verify the environment variables are as expected (sort the slices).
	expectedEnv := getExpectedEnvVars(t, environment)
	cmdEnv := cmd.Env

	sort.Strings(expectedEnv)
	sort.Strings(cmdEnv)

	require.Equal(t, expectedEnv, cmdEnv)
}
