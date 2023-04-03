package limits

import (
	"time"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func (cfg *Config) linuxResources() *specs.LinuxResources {

	lr := specs.LinuxResources{}

	// Set CPU limit, if set.
	if cfg.CPUPercentage != 1.0 {

		// We want to set total CPU time limit. We'll use one year as the period.
		period := uint64(time.Second.Microseconds())
		quota := int64(float64(period) * float64(cfg.CPUPercentage))

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
