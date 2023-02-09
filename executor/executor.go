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

	// Verify the runtime path is valid.
	path := filepath.Join(runtimedir, blocklessCli)
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("invalid runtime path, cli not found (path: %s): %w", path, err)
	}

	e := Executor{
		log:        log,
		workdir:    workdir,
		runtimedir: runtimedir,
	}

	return &e, nil
}
