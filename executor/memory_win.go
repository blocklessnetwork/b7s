//go:build windows
// +build windows

package executor

import (
	"os"
)

// getMemUsage returns process max memory usage in kilobytes.
func getMemUsage(ps *os.ProcessState) int64 {
	// FIXME: See how to retrieve memory usage on windows.
	return 0
}
