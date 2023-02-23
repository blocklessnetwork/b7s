package executor

import (
	"github.com/spf13/afero"
)

// defaultConfig used to create Executor.
var defaultConfig = Config{
	WorkDir:    "workspace",
	RuntimeDir: "",
	FS:         afero.NewOsFs(),
}

// Config represents the Executor configuration.
type Config struct {
	WorkDir    string
	RuntimeDir string
	FS         afero.Fs
}

type Option func(*Config)

// WithWorkDir sets the workspace directory for the executor.
func WithWorkDir(dir string) Option {
	return func(cfg *Config) {
		cfg.WorkDir = dir
	}
}

// WithRuntimeDir sets the runtime directory for the executor.
func WithRuntimeDir(dir string) Option {
	return func(cfg *Config) {
		cfg.RuntimeDir = dir
	}
}

// WithFS sets the FS handler used by the executor.
func WithFS(fs afero.Fs) Option {
	return func(cfg *Config) {
		cfg.FS = fs
	}
}
