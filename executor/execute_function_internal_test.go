package executor

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestExecute_WriteManifest(t *testing.T) {

	const (
		manifestPath = "/manifest.json"
		entry        = "entry-path-value"
		fsRoot       = "fs-root-path-value"
	)

	const expectedManifest = `{
	"fs_root_path": "fs-root-path-value",
	"entry": "entry-path-value",
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
