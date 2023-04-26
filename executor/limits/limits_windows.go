//go:build windows
// +build windows

package limits

import (
	"fmt"

	"golang.org/x/sys/windows"

	"github.com/blocklessnetworking/b7s/models/execute"
)

type Limits struct {
	cfg Config

	jh windows.Handle
}

// New creates a new process resource limit with the given configuration.
func New(opts ...Option) (*Limits, error) {

	// Create job object to which executions will be assigned to.
	name, err := windows.UTF16PtrFromString(DefaultJobObjectName)
	if err != nil {
		return nil, fmt.Errorf("could not prepare job object name: %w", err)
	}

	h, err := windows.CreateJobObject(nil, name)
	if err != nil {
		return nil, fmt.Errorf("could not create job object: %w", err)
	}

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.CPUPercentage < 1.0 {
		err := setCPULimit(h, cfg.CPUPercentage)
		if err != nil {
			windows.CloseHandle(h)
			return nil, fmt.Errorf("could not set CPU limit: %w", err)
		}
	}

	if cfg.MemoryKB > 0 {
		err := setMemLimit(h, cfg.MemoryKB)
		if err != nil {
			windows.CloseHandle(h)
			return nil, fmt.Errorf("could not set memory limit: %w", err)
		}
	}

	l := Limits{
		cfg: cfg,
		jh:  h,
	}

	return &l, nil
}

// LimitProcess will set the resource limits for the process identified by the handle.
func (l *Limits) LimitProcess(proc execute.ProcessID) error {

	err := windows.AssignProcessToJobObject(l.jh, proc.Handle)
	if err != nil {
		return fmt.Errorf("could not assign job to job object: %w", err)
	}

	return nil
}

func (l *Limits) ListProcesses() ([]int, error) {

	pids, err := getJobObjectPids(l.jh)
	if err != nil {
		return nil, fmt.Errorf("could not get processes assigned to job object: %w", err)
	}

	return pids, nil
}

// Close will close the limiter.
func (l *Limits) Close() error {
	return windows.CloseHandle(l.jh)
}
