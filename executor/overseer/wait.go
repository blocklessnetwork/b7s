package overseer

import (
	"errors"
	"fmt"
	"time"
)

func (o *Overseer) Wait(id string) (JobState, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return JobState{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	err := h.cmd.Wait()
	if err != nil {
		return JobState{}, fmt.Errorf("could not wait on job: %w", err)
	}

	endTime := time.Now()

	o.log.Info().Str("stdout", h.stdout.String()).Msg("### observer read stdout")

	state := JobState{
		Status:       StatusDone,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	exitCode := h.cmd.ProcessState.ExitCode()
	state.ExitCode = &exitCode
	if *state.ExitCode != 0 {
		state.Status = StatusFailed
	}

	return state, nil
}
