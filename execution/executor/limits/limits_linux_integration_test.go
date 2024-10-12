//go:build limits && linux
// +build limits,linux

package limits_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/execution/executor/limits"
	"github.com/blocklessnetwork/b7s/models/execute"
)

const (
	cpuMaxFile = "cpu.max"
	memMaxFile = "memory.max"
	pidFile    = "cgroup.procs"
)

func TestLimits(t *testing.T) {

	const (
		cgroup   = limits.DefaultCgroup
		cpuLimit = 0.95
		memLimit = 999_424 // ~1GB rounded to typical page size - 4k
	)

	limiter, err := limits.New(
		limits.WithCgroup(cgroup),
		limits.WithCPUPercentage(cpuLimit),
		limits.WithMemoryKB(memLimit),
	)
	require.NoError(t, err)

	// Always remove all resource limits on end of test.
	defer func() {
		err = limiter.Shutdown()
		require.NoError(t, err)
	}()

	verifyCPULImit(t, cgroup, cpuLimit)
	verifyMemLimit(t, cgroup, memLimit)

	// Verify list of limited processes is empty.
	pids, err := limiter.ListProcesses()
	require.NoError(t, err)
	require.Empty(t, pids)

	// Put resource limit on self.
	// This is effectively a limit on go test so we're conservative with limits.
	proc := execute.ProcessID{
		PID: os.Getpid(),
	}
	err = limiter.LimitProcess(proc)
	require.NoError(t, err)

	// Verify list of limited processes now has a single process.
	pids, err = limiter.ListProcesses()
	require.NoError(t, err)
	require.Len(t, pids, 1)
	require.Equal(t, pids[0], proc.PID)

	// Manually verify the PID limit.
	verifyPids(t, cgroup, []int{proc.PID})
}

func verifyCPULImit(t *testing.T, cgroup string, limit float64) {

	path := filepath.Join(limits.DefaultMountpoint, cgroup, cpuMaxFile)

	payload, err := os.ReadFile(path)
	require.NoError(t, err)

	fields := strings.Fields(string(payload))
	require.Len(t, fields, 2)

	cap, err := strconv.ParseFloat(fields[0], 64)
	require.NoError(t, err)

	period, err := strconv.ParseFloat(fields[1], 64)
	require.NoError(t, err)

	quota := cap / period
	require.Equal(t, limit, quota)
}

func verifyMemLimit(t *testing.T, cgroup string, limitKB int64) {

	path := filepath.Join(limits.DefaultMountpoint, cgroup, memMaxFile)

	payload, err := os.ReadFile(path)
	require.NoError(t, err)

	read := strings.TrimSpace(string(payload))

	limitBytes := limitKB * 1000
	expected := fmt.Sprint(limitBytes)

	require.Equal(t, expected, read)
}

func verifyPids(t *testing.T, cgroup string, pids []int) {
	path := filepath.Join(limits.DefaultMountpoint, cgroup, pidFile)

	payload, err := os.ReadFile(path)
	require.NoError(t, err)

	lines := strings.Split(string(payload), "\n")

	var readPids []int
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Handle trailing newline
		if line == "" {
			continue
		}

		pid, err := strconv.ParseInt(line, 10, 32)
		require.NoError(t, err)

		readPids = append(readPids, int(pid))
	}

	for i, pid := range readPids {
		require.Equal(t, pid, pids[i])
	}
}
