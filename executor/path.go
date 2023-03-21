package executor

import (
	"path/filepath"
)

// requestPaths defines a number of path components relevant to a request.
type requestPaths struct {
	workdir  string
	fsRoot   string
	manifest string
	entry    string
}

func (e *Executor) generateRequestPaths(requestID string, functionID string, method string) requestPaths {

	// Workdir Should be the root for all other paths.
	workdir := filepath.Join(e.cfg.WorkDir, "t", requestID)
	paths := requestPaths{
		workdir:  workdir,
		fsRoot:   filepath.Join(workdir, "fs"),
		manifest: filepath.Join(workdir, "runtime-manifest.json"),
		entry:    filepath.Join(e.cfg.WorkDir, functionID, method),
	}

	return paths
}
