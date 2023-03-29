package limits

import (
	"errors"
	"fmt"

	"github.com/containerd/cgroups"
)

// TODO: Perhaps we crashed and we didn't manage to clean up the cgroup from a previous run.
// Add a cfg flag to override existing cgroup (and overwrite it).

// TODO: Potentially update the cgroup.

// TODO: For now Linux is fine, but try to think cross-platform when it comes to naming, comments etc.

type Limits struct {
	cfg Config

	cgroup cgroups.Cgroup
}

// New creates a new process resource limit with the given configuration.
func New(opts ...Option) (*Limits, error) {

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	// Cgroup should not exist right now.
	_, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(cfg.Cgroup))
	if err == nil {
		return nil, errors.New("cgroup already exists - is there another node instance running?")
	}

	specs := cfg.linuxResources()
	cg, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cfg.Cgroup), specs)
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
func (l *Limits) LimitProcess(pid int) error {

	proc := cgroups.Process{
		Pid: pid,
	}

	err := l.cgroup.Add(proc)
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
