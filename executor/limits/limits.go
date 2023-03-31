package limits

import (
	"errors"
	"fmt"

	"github.com/containerd/cgroups"
	"github.com/containerd/cgroups/v3/cgroup2"
)

// TODO: Perhaps we crashed and we didn't manage to clean up the cgroup from a previous run.
// Add a cfg flag to override existing cgroup (and overwrite it).

// TODO: Potentially update the cgroup.

// TODO: For now Linux is fine, but try to think cross-platform when it comes to naming, comments etc.

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

	// Cgroup should not exist right now.
	// _, err := cgroup2.LoadSystemd(cfg.Cgroup, cfg.Cgroup)
	// if err == nil {
	// 	return nil, errors.New("cgroup already exists - is there another node instance running?")
	// }

	specs := cfg.cgroupV2Resources()
	cg, err := cgroup2.NewSystemd(cfg.Cgroup, cfg.Cgroup, -1, specs)
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
func (l *Limits) LimitProcess(pid uint64) error {

	err := l.cgroup.AddProc(pid)
	if err != nil {
		return fmt.Errorf("could not set resouce limit for the process: %w", err)
	}

	return nil
}

// Remove removes the created resource limit.
func (l *Limits) Remove() error {
	err := l.cgroup.Delete()
	if err != nil {
		return fmt.Errorf("could not remove resource limits: %w", err)
	}

	return nil
}
