package limits

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Cgroup(t *testing.T) {

	const cgroup = "/blockless-test"

	cfg := Config{
		Cgroup: DefaultCgroup,
	}

	WithCgroup(cgroup)(&cfg)
	require.Equal(t, cgroup, cfg.Cgroup)
}

func TestConfig_WithCPUPercentage(t *testing.T) {

	const pct = 0.7

	cfg := Config{
		CPUPercentage: 1.0,
	}

	WithCPUPercentage(pct)(&cfg)
	require.Equal(t, pct, cfg.CPUPercentage)
}

func TestConfig_WithMemoryKB(t *testing.T) {

	const limit = int64(200_000)

	cfg := Config{
		MemoryKB: 10,
	}

	WithMemoryKB(limit)(&cfg)
	require.Equal(t, limit, cfg.MemoryKB)
}

func TestConfig_JobName(t *testing.T) {

	const jobName = "blockless-test"

	cfg := Config{
		JobName: DefaultJobObjectName,
	}

	WithJobObjectName(jobName)(&cfg)
	require.Equal(t, jobName, cfg.JobName)
}
