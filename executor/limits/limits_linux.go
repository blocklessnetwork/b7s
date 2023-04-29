//go:build linux
// +build linux

package limits

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/containerd/cgroups/v3"
	"github.com/containerd/cgroups/v3/cgroup2"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// TODO: Add support for cgroups v1 - determine on the fly which version to use

type Limits struct {
	cfg Config

	cgroup *cgroup2.Manager
}

// New creates a new process resource limit with the given configuration.
func New(opts ...Option) (*Limits, error) {

	// Check if the system supports cgroups v2.
	haveV2 := false
	if cgroups.Mode() == cgroups.Unified {
		haveV2 = true
	}
	if !haveV2 {
		return nil, errors.New("cgroups v2 is not supported")
	}

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	specs := cfg.cgroupV2Resources()

	// NOTE: Library we use for handling cgroups will also remove the directory on failure.
	// Since we need root privileges to create it, this can cause problems.
	cg, err := cgroup2.NewManager(DefaultMountpoint, cfg.Cgroup, specs)
	if err != nil {
		return nil, fmt.Errorf("could not create cgroup: %w", err)
	}

	l := Limits{
		cfg:    cfg,
		cgroup: cg,
	}

	return &l, nil
}

// LimitProcess will set the resource limits for the process with the given PID.
func (l *Limits) LimitProcess(proc execute.ProcessID) error {

	pid := proc.PID
	err := l.cgroup.AddProc(uint64(pid))
	if err != nil {
		return fmt.Errorf("could not set resouce limit for process (pid: %v): %w", pid, err)
	}

	return nil
}

// ListProcesses will return the pids of the processes that were added to the resource limit group.
func (l *Limits) ListProcesses() ([]int, error) {

	var list []int
	pids, err := l.cgroup.Procs(false)
	if err != nil {
		return nil, fmt.Errorf("could not get list of limited processes: %w", err)
	}

	for _, pid := range pids {
		list = append(list, int(pid))
	}

	return list, nil
}

// Shutdown will remove any set resource limits.
func (l *Limits) Shutdown() error {

	// Remove all limits effectively sets them to very large values, which is different from "removing" them.
	period := uint64(time.Second.Microseconds())
	memLimit := int64(math.MaxInt64)

	resources := cgroup2.Resources{
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(nil, &period),
		},
		Memory: &cgroup2.Memory{
			Max: &memLimit,
		},
	}

	err := l.cgroup.Update(&resources)
	if err != nil {
		return fmt.Errorf("could not update resource limits: %v", err)
	}

	return nil
}
