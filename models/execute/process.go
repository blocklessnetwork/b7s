package execute

// ProcessID is used to identify an OS process.
type ProcessID struct {
	PID    int     // PID can used to identify a process on all platforms.
	Handle uintptr // windows.Handle value that can be used for Windows-specific operations.
}
