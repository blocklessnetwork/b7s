package function

import (
	"path/filepath"
	"strings"
)

// cleanPath will return the path relative to the workdir.
func (h *Handler) cleanPath(path string) string {
	return filepath.Clean(strings.TrimPrefix(path, h.workdir))
}
