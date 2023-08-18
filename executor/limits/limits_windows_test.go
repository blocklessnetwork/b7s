//go:build windows
// +build windows

package limits_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"

	"github.com/blocklessnetwork/b7s/executor/limits"
	"github.com/blocklessnetwork/b7s/models/execute"
)

func TestLimits(t *testing.T) {

	// NOTE: We currently don't have an implementation to query set resource limits.
	// What we can do is spin up another process and limit it by assigning it to the job object.
	// We can then query if the process PID is indeed assigned - verifying that the job object
	// exists and the process has been assigned to it, providing *some* assurance.

	const (
		cpuLimit = 0.95
		memLimit = 256_000

		// We use ping to simulate "sleep", more info below.
		pingPath = `C:\Windows\System32\ping.exe`
	)

	limiter, err := limits.New(
		limits.WithCPUPercentage(cpuLimit),
		limits.WithMemoryKB(memLimit),
	)
	require.NoError(t, err)

	defer func() {
		err = limiter.Shutdown()
		require.NoError(t, err)
	}()

	// Verify list of limited processes is empty.
	pids, err := limiter.ListProcesses()
	require.NoError(t, err)
	require.Empty(t, pids)

	// Ideally we want to run something like "sleep" for a few seconds.
	// However, "sleep" is not present by default on Windows. There is a "timeout" alternative,
	// but it does not support running programmatically because of input redirection - it expects
	// user input to stop waiting. Instead we'll just use a "ping" on the loopback address.
	// => See https://www.ibm.com/support/pages/timeout-command-run-batch-job-exits-immediately-and-returns-error-input-redirection-not-supported-exiting-process-immediately
	cmd := exec.Command(pingPath, "-n", "1", "127.0.0.1")

	err = cmd.Start()
	require.NoError(t, err)

	handle, err := windows.OpenProcess(
		windows.PROCESS_TERMINATE|windows.PROCESS_SET_QUOTA,
		false,
		uint32(cmd.Process.Pid),
	)
	require.NoError(t, err)

	defer windows.CloseHandle(handle)

	proc := execute.ProcessID{
		PID:    cmd.Process.Pid,
		Handle: uintptr(handle),
	}
	err = limiter.LimitProcess(proc)
	require.NoError(t, err)

	pids, err = limiter.ListProcesses()
	require.NoError(t, err)

	require.Len(t, pids, 1)
	require.Equal(t, cmd.Process.Pid, pids[0])

	err = cmd.Wait()
	require.NoError(t, err)
}
