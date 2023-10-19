package overseer

import (
	"errors"
	"time"
)

func (o *Overseer) Observe(id string) (any, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return nil, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	state := JobState{
		Status:       StatusRunning,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		ObservedTime: time.Now(),
	}

	if h.cmd.ProcessState != nil {
		state.Status = StatusDone

		// TODO: Wait on process in a goroutine to have accurate end time.
		endTime := time.Now()
		state.EndTime = &endTime

		exitCode := h.cmd.ProcessState.ExitCode()
		state.ExitCode = &exitCode
		if *state.ExitCode != 0 {
			state.Status = StatusFailed
		}
	}

	return state, nil
}
