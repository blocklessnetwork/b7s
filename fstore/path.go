package fstore

import (
	"path/filepath"
	"strings"
)

// cleanPath will return the path relative to the workdir.
func (f *FStore) cleanPath(path string) string {
	return filepath.Clean(strings.TrimPrefix(path, f.workdir))
}
