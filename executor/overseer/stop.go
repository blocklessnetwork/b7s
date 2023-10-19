package overseer

import (
	"errors"
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

	h.cancel()

	endTime := time.Now()

	state := JobState{
		Status:       StatusKilled,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	if h.cmd.ProcessState != nil {
		exitCode := h.cmd.ProcessState.ExitCode()
		state.ExitCode = &exitCode
		if *state.ExitCode != 0 {
			state.Status = StatusFailed
		}
	}

	return state, nil
}
