package overseer

import (
	"errors"
	"os"
	"syscall"
)

func determineProcessStatus(ps *os.ProcessState) (*int, JobStatus) {

	if ps == nil {
		return nil, StatusRunning
	}

	exitCode := ps.ExitCode()
	if exitCode == 0 {
		return &exitCode, StatusDone
	}

	signaled, _ := wasSignalled(ps)
	if signaled {
		return &exitCode, StatusKilled
	}

	return &exitCode, StatusFailed
}

// TODO: Check - OS dependent.
func wasSignalled(ps *os.ProcessState) (bool, error) {

	ws, ok := ps.Sys().(syscall.WaitStatus)
	if !ok {
		return false, errors.New("unexpected type for exit information")
	}

	return ws.Signaled(), nil
}
