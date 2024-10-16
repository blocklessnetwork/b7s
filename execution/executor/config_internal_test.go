package executor

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestWithWorkDir(t *testing.T) {

	const workdir = "/var/tmp/b7s"

	cfg := Config{
		WorkDir: "",
	}

	WithWorkDir(workdir)(&cfg)
	require.Equal(t, workdir, cfg.WorkDir)
}

func TestWithRuntimePath(t *testing.T) {

	const runtimePath = "/usr/local/bin/blockless-cli"

	cfg := Config{
		RuntimePath: "",
	}

	WithRuntimePath(runtimePath)(&cfg)
	require.Equal(t, runtimePath, cfg.RuntimePath)
}

func TestWithFS(t *testing.T) {

	var fs = afero.NewOsFs()

	cfg := Config{
		FS: afero.NewMemMapFs(),
	}

	WithFS(fs)(&cfg)
	require.Equal(t, fs, cfg.FS)
}
