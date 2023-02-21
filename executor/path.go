package executor

import (
	"path"
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
	workdir := path.Join(e.workdir, "t", requestID)
	paths := requestPaths{
		workdir:  workdir,
		fsRoot:   path.Join(workdir, "fs"),
		manifest: path.Join(workdir, "runtime-manifest.json"),
		entry:    path.Join(e.workdir, functionID, method),
	}

	return paths
}
