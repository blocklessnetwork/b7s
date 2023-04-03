package executor

import (
	"github.com/spf13/afero"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// defaultConfig used to create Executor.
var defaultConfig = Config{
	WorkDir:        "workspace",
	RuntimeDir:     "",
	ExecutableName: blockless.RuntimeCLI(),
	FS:             afero.NewOsFs(),
}

// Config represents the Executor configuration.
type Config struct {
	WorkDir        string   // directory where files needed for the execution are stored
	RuntimeDir     string   // directory where the executable can be found
	ExecutableName string   // name for the executable
	FS             afero.Fs // FS accessor
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

// WithExecutableName sets the name of the executable that should be ran.
func WithExecutableName(name string) Option {
	return func(cfg *Config) {
		cfg.ExecutableName = name
	}
}
