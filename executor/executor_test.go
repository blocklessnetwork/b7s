package executor_test

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestExecutor_Create(t *testing.T) {
	t.Run("nominal case", func(t *testing.T) {

		var (
			runtimeDir = os.TempDir()
			cliPath    = path.Join(runtimeDir, blockless.RuntimeCLI())
			fs         = afero.NewMemMapFs()
		)

		_, err := fs.Create(cliPath)
		require.NoError(t, err)

		_, err = executor.New(mocks.NoopLogger,
			executor.WithRuntimeDir(runtimeDir),
			executor.WithFS(fs),
		)
		require.NoError(t, err)
	})
	t.Run("missing runtime path", func(t *testing.T) {

		const (
			runtimeDir = "/usr/local/bin"
		)

		// Use empty FS that surely will not have the runtime anywhere.
		executor, err := executor.New(mocks.NoopLogger,
			executor.WithRuntimeDir(runtimeDir),
			executor.WithFS(afero.NewMemMapFs()),
		)
		require.Error(t, err)
		require.Nil(t, executor)
	})

}
