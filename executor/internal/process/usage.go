package process

import (
	"fmt"
	"os/exec"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// GetUsage returns the resource usage information about the executed process.
func GetUsage(cmd *exec.Cmd) (execute.Usage, error) {

	// If `cmd.Process` is empty, it means that the command has not been executed yet.
	// On the other hand, if the `ProcessState` is empty - it means that the process did not yet
	// complete, or was not Wait-ed on.
	if cmd.Process == nil || cmd.ProcessState == nil {
		return execute.Usage{}, fmt.Errorf("process not started or not yet completed")
	}

	ps := cmd.ProcessState

	usage := execute.Usage{
		CPUUserTime: ps.UserTime(),
		CPUSysTime:  ps.SystemTime(),
		MemoryMaxKB: getMemUsage(ps),
	}

	return usage, nil
}
