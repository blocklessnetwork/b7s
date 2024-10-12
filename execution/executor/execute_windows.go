//go:build windows
// +build windows

package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/sys/windows"

	"github.com/blocklessnetwork/b7s/executor/internal/process"
	"github.com/blocklessnetwork/b7s/models/execute"
)

// executeCommand on Windows contains some dark sorcery. On Windows, the `rusage` equivalent does not include
// memory information. In order to get this info, we need the process `handle`, not just its PID. Process
// handle can be obtained by using `OpenProcess` syscall, but that is a data race, as the process might have
// already exited by the time our syscall returns. To do this, we rely on the fact that the stdlib does not
// change the process handle until a successful `Wait`. And on Windows, as long as we hold the handle, we
// have access to the process information. So we'll use reflection to get the value of the handle and do a
// `DuplicateHandle“ syscall. With this duplicated handle, we'll be able to access all the info we need.
// Additionally, the `DuplicateHandle` syscall will fail if we do anything wrong, so it will also act as a
// validation layer.
func (e *Executor) executeCommand(cmd *exec.Cmd) (execute.RuntimeOutput, execute.Usage, error) {

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command and collect output.
	start := time.Now()
	err := cmd.Start()
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not start process: %w", err)
	}

	childHandle, err := process.ReadHandle(cmd)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not get process handle: %w", err)
	}

	// Create a duplicate handle - only for me (current process), not inheritable.
	var handle windows.Handle
	me := windows.CurrentProcess()
	err = windows.DuplicateHandle(
		me,
		childHandle,
		me,
		&handle,
		windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_TERMINATE|windows.PROCESS_SET_QUOTA,
		false,
		0)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not duplicate process handle: %w", err)
	}
	defer func() {
		err := windows.CloseHandle(handle)
		if err != nil {
			e.log.Error().Err(err).Int("pid", cmd.Process.Pid).Msg("could not close handle")
		}
	}()

	proc := execute.ProcessID{
		PID:    cmd.Process.Pid,
		Handle: uintptr(handle),
	}
	err = e.cfg.Limiter.LimitProcess(proc)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not set resource limits: %w", err)
	}

	// Now we can safely wait for the child process to complete.
	cmdErr := cmd.Wait()
	end := time.Now()

	out := execute.RuntimeOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	// Create usage information.
	duration := end.Sub(start)
	usage, err := process.GetUsage(cmd)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not retrieve usage data: %w", err)
	}

	// Returned memory usage is in bytes, so convert it to kilobytes.
	mem, err := process.GetMemUsageForHandle(handle)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not retrieve memory data: %w", err)
	}

	usage.MemoryMaxKB = int64(mem) / 1000
	usage.WallClockTime = duration

	if cmdErr != nil {
		return out, usage, fmt.Errorf("process execution failed: %w", cmdErr)
	}

	return out, usage, nil
}
