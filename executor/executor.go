package executor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

// Executor provides the capabilities to run external applications.
type Executor struct {
	log zerolog.Logger

	workdir    string
	runtimedir string
}

// New creates a new Executor with the specified working directory.
func New(log zerolog.Logger, workdir string, runtimedir string) (*Executor, error) {

	// We need the absolute path for the runtime, since we'll be changing
	// the working directory on execution.
	runtime, err := filepath.Abs(runtimedir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for runtime (path: %s): %w", runtimedir, err)
	}

	// Verify the runtime path is valid.
	cliPath := filepath.Join(runtime, blocklessCli)
	_, err = os.Stat(cliPath)
	if err != nil {
		return nil, fmt.Errorf("invalid runtime path, cli not found (path: %s): %w", cliPath, err)
	}

	e := Executor{
		log:        log,
		workdir:    workdir,
		runtimedir: runtime,
	}

	return &e, nil
}
