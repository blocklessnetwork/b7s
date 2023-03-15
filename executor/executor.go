package executor

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog"
)

// Executor provides the capabilities to run external applications.
type Executor struct {
	log zerolog.Logger
	cfg Config
}

// New creates a new Executor with the specified working directory.
func New(log zerolog.Logger, options ...Option) (*Executor, error) {

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	// We need the absolute path for the runtime, since we'll be changing
	// the working directory on execution.
	runtime, err := filepath.Abs(cfg.RuntimeDir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for runtime (path: %s): %w", cfg.RuntimeDir, err)
	}
	cfg.RuntimeDir = runtime

	// Verify the runtime path is valid.
	cliPath := filepath.Join(cfg.RuntimeDir, cfg.ExecutableName)
	_, err = cfg.FS.Stat(cliPath)
	if err != nil {
		return nil, fmt.Errorf("invalid runtime path, cli not found (path: %s): %w", cliPath, err)
	}

	e := Executor{
		log: log.With().Str("component", "executor").Logger(),
		cfg: cfg,
	}

	return &e, nil
}
