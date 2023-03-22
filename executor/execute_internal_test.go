package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestExecute_CreateCMD(t *testing.T) {

	var (
		runtimeDir     = "/usr/local/bin"
		workdir        = "/var/tmp/b7s"
		functionID     = "function-id"
		functionMethod = "function-method"

		executablePath = filepath.Join(runtimeDir, blocklessCli)

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
			ExecutableName: blocklessCli,
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

func TestExecute_WriteManifest(t *testing.T) {

	const (
		manifestPath = "/manifest.json"
		entry        = "entry-path-value"
		fsRoot       = "fs-root-path-value"
	)

	const expectedManifest = `{
	"fs_root_path": "fs-root-path-value",
	"entry": "entry-path-value",
	"limited_fuel": 100000000,
	"limited_memory": 200,
	"permissions": [
		"permission-string-a",
		"permission-string-b",
		"permission-string-c"
	]
}`

	request := execute.Request{
		Config: execute.Config{
			Permissions: []string{
				"permission-string-a",
				"permission-string-b",
				"permission-string-c",
			},
		},
	}
	paths := requestPaths{
		manifest: manifestPath,
		entry:    entry,
		fsRoot:   fsRoot,
	}

	fs := afero.NewMemMapFs()
	executor := Executor{
		log: mocks.NoopLogger,
		cfg: Config{
			FS: fs,
		},
	}

	err := executor.writeExecutionManifest(request, paths)
	require.NoError(t, err)

	read, err := afero.ReadFile(fs, manifestPath)
	require.NoError(t, err)
	require.Equal(t, expectedManifest, string(read))
}

func TestExecute_WriteFile(t *testing.T) {

	const (
		filename = "dummy-file.txt"
		content  = "This is the content of the file we want to create."
	)
	fs := afero.NewMemMapFs()

	executor := Executor{
		log: mocks.NoopLogger,
		cfg: Config{
			FS: fs,
		},
	}

	err := executor.writeFile(filename, []byte(content))
	require.NoError(t, err)

	read, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	require.Equal(t, content, string(read))
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

	list := fmt.Sprintf("%s=%s", blsListEnvName, strings.Join(names, ";"))
	out = append(out, list)

	return out
}
