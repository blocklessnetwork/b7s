//go:build !windows
// +build !windows

package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/blocklessnetworking/b7s/executor/internal/process"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// executeCommand on non-windows systems is pretty straightforward and equivalent to the ordinary `cmd.Run()` or `cmd.Output`.
func (e *Executor) executeCommand(cmd *exec.Cmd) (string, execute.Usage, error) {

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Execute the command and collect output.
	start := time.Now()
	err := cmd.Start()
	if err != nil {
		return "", execute.Usage{}, fmt.Errorf("could not start process: %w", err)
	}

	// Set resource limits on the process.
	e.log.Debug().Int("pid", cmd.Process.Pid).Msg("setting resource limits for process")

	err = e.cfg.Limiter.LimitProcess(uint64(cmd.Process.Pid))
	if err != nil {
		return "", execute.Usage{}, fmt.Errorf("could not limit process: %w", err)
	}

	e.log.Debug().Int("pid", cmd.Process.Pid).Msg("resource limits set for process")

	err = cmd.Wait()
	if err != nil {
		return "", execute.Usage{}, fmt.Errorf("could not wait on process: %w", err)
	}

	end := time.Now()

	// Create usage information.
	duration := end.Sub(start)
	usage, err := process.GetUsage(cmd)
	if err != nil {
		return "", execute.Usage{}, fmt.Errorf("could not retrieve usage data: %w", err)
	}

	usage.WallClockTime = duration

	return stdout.String(), usage, nil
}
