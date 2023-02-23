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

func TestWithRuntimeDir(t *testing.T) {

	const runtimeDir = "/usr/local/bin"

	cfg := Config{
		RuntimeDir: "",
	}

	WithRuntimeDir(runtimeDir)(&cfg)
	require.Equal(t, runtimeDir, cfg.RuntimeDir)
}

func TestWithFS(t *testing.T) {

	var fs = afero.NewOsFs()

	cfg := Config{
		FS: afero.NewMemMapFs(),
	}

	WithFS(fs)(&cfg)
	require.Equal(t, fs, cfg.FS)
}
