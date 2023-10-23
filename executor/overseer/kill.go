package overseer

import (
	"errors"
	"fmt"
	"time"
)

func (o *Overseer) Kill(id string) (JobState, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return JobState{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	if h.cmd.Process == nil {
		return JobState{}, errors.New("job is not running")
	}

	// TODO: Check if the process has started any processes.
	err := h.cmd.Process.Kill()
	if err != nil {
		return JobState{}, fmt.Errorf("could not kill process: %w", err)
	}

	endTime := time.Now()

	// Until we wait on the process (read its state) it will still exist on the system as a zombie process.
	err = h.cmd.Wait()
	if err != nil {

		signaled, _ := wasSignalled(h.cmd.ProcessState)
		if !signaled {
			return JobState{}, fmt.Errorf("could not wait on process: %w", err)
		}

		o.log.Trace().Err(err).Str("job", id).Msg("expected error - job was signaled, wait produced an error")
	}

	state := JobState{
		Status:       StatusKilled,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	exitCode, status := determineProcessStatus(h.cmd.ProcessState)
	state.ExitCode = exitCode
	state.Status = status

	return state, nil
}
