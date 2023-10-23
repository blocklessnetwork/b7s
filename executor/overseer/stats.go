package overseer

import (
	"errors"
	"time"
)

func (o *Overseer) Stats(id string) (JobState, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return JobState{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	state := JobState{
		Status:       StatusRunning,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		ObservedTime: time.Now(),
		// TODO: Process stats.
	}

	exitCode, status := determineProcessStatus(h.cmd.ProcessState)
	state.ExitCode = exitCode
	state.Status = status

	if h.cmd.ProcessState != nil {
		// NOTE: Perhaps wait on process in a goroutine to have accurate end time.
		endTime := time.Now()
		state.EndTime = &endTime
	}

	return state, nil
}
