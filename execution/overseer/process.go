package overseer

import (
	"errors"
	"os"
	"syscall"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

func determineProcessStatus(ps *os.ProcessState) (*int, job.Status) {

	if ps == nil {
		return nil, job.StatusRunning
	}

	exitCode := ps.ExitCode()
	if exitCode == 0 {
		return &exitCode, job.StatusDone
	}

	signaled, _ := wasSignalled(ps)
	if signaled {
		return &exitCode, job.StatusKilled
	}

	return &exitCode, job.StatusFailed
}

// TODO: Check - OS dependent.
func wasSignalled(ps *os.ProcessState) (bool, error) {

	ws, ok := ps.Sys().(syscall.WaitStatus)
	if !ok {
		return false, errors.New("unexpected type for exit information")
	}

	return ws.Signaled(), nil
}
