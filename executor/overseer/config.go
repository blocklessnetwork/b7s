package overseer

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

var defaultConfig = Config{
	Workdir: "workspace",
	FS:      afero.NewOsFs(),
}

type Config struct {
	Workdir   string
	FS        afero.Fs
	Allowlist []string
}

type Option func(*Config)

// WithWorkdir sets the workspace directory for the overseer.
func WithWorkdir(dir string) Option {
	return func(cfg *Config) {
		cfg.Workdir = dir
	}
}

// WithFS sets the FS handler used by the overseer.
func WithFS(fs afero.Fs) Option {
	return func(cfg *Config) {
		cfg.FS = fs
	}
}

// WithAllowlist specifies the executables the overseer is allowed to start.
// Executables should be listed using their full paths.
func WithAllowlist(executables []string) Option {
	return func(cfg *Config) {
		cfg.Allowlist = executables
	}
}

func (cfg Config) Validate() error {

	for _, path := range cfg.Allowlist {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("path from allowlist not absolute: %s", path)
		}
	}

	return nil
}
