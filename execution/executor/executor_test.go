package executor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/execution/executor"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestExecutor_Create(t *testing.T) {
	t.Run("nominal case", func(t *testing.T) {

		var (
			cliPath = filepath.Join(os.TempDir(), blockless.RuntimeCLI())
			fs      = afero.NewMemMapFs()
		)

		_, err := fs.Create(cliPath)
		require.NoError(t, err)

		_, err = executor.New(mocks.NoopLogger,
			executor.WithRuntimePath(cliPath),
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
			executor.WithRuntimePath(runtimeDir),
			executor.WithFS(afero.NewMemMapFs()),
		)
		require.Error(t, err)
		require.Nil(t, executor)
	})

}
