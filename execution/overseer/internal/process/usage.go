package process

import (
	"fmt"
	"os/exec"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

func GetUsage(cmd *exec.Cmd) (job.ResourceUsage, error) {

	// If `cmd.Process` is empty, it means that the command has not been executed yet.
	// On the other hand, if the `ProcessState` is empty - it means that the process did not yet
	// complete, or was not Wait-ed on.
	if cmd.Process == nil || cmd.ProcessState == nil {
		return job.ResourceUsage{}, fmt.Errorf("process not started or not yet completed")
	}

	ps := cmd.ProcessState

	usage := job.ResourceUsage{
		CPUUserTime: ps.UserTime(),
		CPUSysTime:  ps.SystemTime(),
		MemoryMaxKB: getMemUsage(ps),
	}

	return usage, nil
}
