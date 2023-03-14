//go:build !windows
// +build !windows

package executor

import (
	"os"
	"syscall"
)

// getMemUsage returns process max memory usage in kilobytes.
func getMemUsage(ps *os.ProcessState) int64 {

	usage, ok := ps.SysUsage().(*syscall.Rusage)
	if !ok {
		return 0
	}

	return usage.Maxrss
}
