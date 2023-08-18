//go:build !linux && !windows
// +build !linux,!windows

package limits

import (
	"errors"

	"github.com/blocklessnetwork/b7s/models/execute"
)

// NOTE: Placeholder for operating systems where we do not support limiters yet.

type Limits struct {
	cfg Config
}

// New creates a new process resource limit with the given configuration.
func New(opts ...Option) (*Limits, error) {
	return nil, errors.New("TBD: not implemented")
}

// LimitProcess will set the resource limits for the process with the given PID.
func (l *Limits) LimitProcess(proc execute.ProcessID) error {
	return errors.New("TBD: not implemented")
}

// ListProcesses will return the pids of the processes that were added to the resource limit group.
func (l *Limits) ListProcesses() ([]int, error) {
	return nil, errors.New("TBD: not implemented")
}

// RemoveAllLimits will remove any set resource limits.
func (l *Limits) RemoveAllLimits() error {
	return errors.New("TBD: not implemented")
}

// Close will close the limiter.
func (l *Limits) Shutdown() error {
	return errors.New("TBD: not implemented")
}
