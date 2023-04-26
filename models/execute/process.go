package execute

import (
	"golang.org/x/sys/windows"
)

// ProcessID is used to identify an OS process.
type ProcessID struct {
	PID    int            // PID can used to identify a process on all platforms.
	Handle windows.Handle // Handle can be used for Windows-specific operations.
}
