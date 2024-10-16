package executor

import (
	"github.com/armon/go-metrics"
	"github.com/spf13/afero"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// defaultConfig used to create Executor.
var defaultConfig = Config{
	WorkDir:         "workspace",
	RuntimePath:     blockless.RuntimeCLI(),
	FS:              afero.NewOsFs(),
	Limiter:         &noopLimiter{},
	DriversRootPath: "",
}

// Config represents the Executor configuration.
type Config struct {
	WorkDir         string           // directory where files needed for the execution are stored
	RuntimePath     string           // full path to the runtime
	DriversRootPath string           // where are cgi drivers stored
	FS              afero.Fs         // FS accessor
	Limiter         Limiter          // Resource limiter for executed processes
	Metrics         *metrics.Metrics // Metrics handle
}

type Option func(*Config)

// WithWorkDir sets the workspace directory for the executor.
func WithWorkDir(dir string) Option {
	return func(cfg *Config) {
		cfg.WorkDir = dir
	}
}

// WithRuntimePath sets the path to the runtime.
func WithRuntimePath(path string) Option {
	return func(cfg *Config) {
		cfg.RuntimePath = path
	}
}

// WithFS sets the FS handler used by the executor.
func WithFS(fs afero.Fs) Option {
	return func(cfg *Config) {
		cfg.FS = fs
	}
}

// WithLimiter sets the resource limiter called for each individual execution.
func WithLimiter(limiter Limiter) Option {
	return func(cfg *Config) {
		cfg.Limiter = limiter
	}
}

// WithMetrics sets the metrics handler.
func WithMetrics(metrics *metrics.Metrics) Option {
	return func(cfg *Config) {
		cfg.Metrics = metrics
	}
}
