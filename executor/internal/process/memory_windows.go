//go:build windows
// +build windows

package process

import (
	"fmt"
	"os"

	"github.com/Microsoft/go-winio/pkg/process"
	"golang.org/x/sys/windows"
)

// GetMemUsageForHandle returns the peak working set size for the process, in bytes.
func GetMemUsageForHandle(handle windows.Handle) (uint, error) {

	counters, err := process.GetProcessMemoryInfo(handle)
	if err != nil {
		return 0, fmt.Errorf("could not get memory info for handle: %s", err)
	}

	return counters.PeakWorkingSetSize, nil
}

// getMemUsage is not implemented on Windows. See `GetMemUsageForHandle`.
func getMemUsage(ps *os.ProcessState) int64 {
	return 0
}
