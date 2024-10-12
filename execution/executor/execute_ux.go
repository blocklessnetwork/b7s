//go:build !windows
// +build !windows

package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/blocklessnetwork/b7s/execution/executor/internal/process"
	"github.com/blocklessnetwork/b7s/models/execute"
)

// executeCommand on non-windows systems is pretty straightforward and equivalent to the ordinary `cmd.Run()` or `cmd.Output`.
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

	proc := execute.ProcessID{
		PID: cmd.Process.Pid,
	}
	err = e.cfg.Limiter.LimitProcess(proc)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not set resource limits: %w", err)
	}

	// Return execution error with as much info below.
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

	usage.WallClockTime = duration

	if cmdErr != nil {
		return out, usage, fmt.Errorf("process execution failed: %w", cmdErr)
	}

	return out, usage, nil
}
