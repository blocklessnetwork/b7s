package limits

import (
	"math"
	"time"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func limitsToResources(limits Limits) *cgroup2.Resources {

	var (
		period                = uint64(time.Second.Microseconds())
		unlimitedMemory int64 = math.MaxInt64
	)

	// By default, remove any previous limits.
	lr := specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Period: &period,
			Quota:  nil,
		},
		Memory: &specs.LinuxMemory{
			Limit: &unlimitedMemory,
		},
		Pids: &specs.LinuxPids{
			Limit: 0,
		},
	}

	if limits.CPUPercentage > 0 && limits.CPUPercentage < 1.0 {
		quota := int64(float64(period) * limits.CPUPercentage)
		lr.CPU.Quota = &quota
	}

	memLimit := limits.MemoryKB * 1000
	if memLimit > 0 {
		lr.Memory.Limit = &memLimit
	}

	if limits.ProcLimit > 0 {
		lr.Pids.Limit = int64(limits.ProcLimit)
	}

	return cgroup2.ToResources(&lr)
}
