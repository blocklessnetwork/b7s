package executor

import (
	"path/filepath"
)

func (e *Executor) generateDirName(requestID string) string {

	dir := filepath.Join(e.workdir, "t", requestID, "fs")
	return dir
}
