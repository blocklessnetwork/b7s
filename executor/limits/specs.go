package limits

import (
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func (cfg *Config) linuxResources() *specs.LinuxResources {

	lr := specs.LinuxResources{}

	if cfg.CPUTime > 0 {

		// We want to set total CPU time limit. We'll use one year as the period.
		period := uint64(year.Microseconds())
		quota := cfg.CPUTime.Microseconds()

		lr.CPU = &specs.LinuxCPU{
			Period: &period,
			Quota:  &quota,
		}
	}

	// Set memory limit, if set.
	if cfg.MemoryKB > 0 {

		// Convert limit to bytes.
		memLimit := cfg.MemoryKB * 1000

		lr.Memory = &specs.LinuxMemory{
			Limit: &memLimit,
		}
	}

	return &lr
}

func (cfg *Config) cgroupV2Resources() *cgroup2.Resources {
	lr := cfg.linuxResources()
	return cgroup2.ToResources(lr)
}
