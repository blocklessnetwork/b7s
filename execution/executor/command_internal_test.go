package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestExecute_CreateCMD(t *testing.T) {

	var (
		runtimeDir     = "/usr/local/bin"
		workdir        = "/var/tmp/b7s"
		functionID     = "function-id"
		functionMethod = "function-method"
		runtimeLogger  = "whatever.log"
		limitedMemory  = 256

		executablePath = filepath.Join(runtimeDir, blockless.RuntimeCLI())

		requestID   = mocks.GenericUUID.String()
		stdin       = "dummy stdin payload"
		environment = getEnvVars(t)

		request = execute.Request{
			Config: execute.Config{
				Stdin:       &stdin,
				Environment: environment,
				Runtime: execute.BLSRuntimeConfig{
					Memory: uint64(limitedMemory),
					Logger: runtimeLogger,
				},
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

	// NOTE: This verification of flags is pretty rigid - it's expects flags in a specific order.

	// Verify CLI arguments are correct.
	require.Len(t, cmd.Args, 9)
	require.Equal(t, executablePath, cmd.Args[0])
	require.Equal(t, paths.input, cmd.Args[1])

	require.Equal(t, "--"+execute.BLSRuntimeFlagFSRoot, cmd.Args[2])
	require.Equal(t, paths.fsRoot, cmd.Args[3])

	require.Equal(t, "--"+execute.BLSRuntimeFlagMemory, cmd.Args[4])
	require.Equal(t, fmt.Sprint(limitedMemory), cmd.Args[5])

	require.Equal(t, "--"+execute.BLSRuntimeFlagLogger, cmd.Args[6])
	require.Equal(t, runtimeLogger, cmd.Args[7])

	require.Equal(t, "--", cmd.Args[8])

	// Verify working directory is correct.
	require.Equal(t, paths.workdir, cmd.Dir)

	// Verify the environment variables are as expected (sort the slices).
	expectedEnv := getExpectedEnvVars(t, environment)
	cmdEnv := cmd.Env

	sort.Strings(expectedEnv)
	sort.Strings(cmdEnv)

	require.Equal(t, expectedEnv, cmdEnv)
}

func getEnvVars(t *testing.T) []execute.EnvVar {
	t.Helper()

	const (
		nameRoot  = "executor-env-var-name"
		valueRoot = "executor-env-var-value"
		count     = 10
	)

	env := make([]execute.EnvVar, 0, count)
	for i := 0; i < count; i++ {

		e := execute.EnvVar{
			Name:  fmt.Sprintf("%s-%d", nameRoot, i),
			Value: fmt.Sprintf("%s-%d", valueRoot, i),
		}

		env = append(env, e)
	}

	return env
}

// getExpectedEnvVars return the complete list of environment variables expected by the CLI.
func getExpectedEnvVars(t *testing.T, environment []execute.EnvVar) []string {
	t.Helper()

	out := os.Environ()

	names := make([]string, 0, len(environment))
	for _, env := range environment {
		e := fmt.Sprintf("%s=%s", env.Name, env.Value)
		out = append(out, e)

		names = append(names, env.Name)
	}

	list := fmt.Sprintf("%s=%s", blockless.RuntimeEnvVarList, strings.Join(names, ";"))
	out = append(out, list)

	return out
}
