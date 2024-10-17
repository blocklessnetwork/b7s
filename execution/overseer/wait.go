package overseer

import (
	"errors"
	"time"

	"github.com/blocklessnetwork/b7s/execution/overseer/internal/process"
	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

func (o *Overseer) Wait(id string) (job.State, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return job.State{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	defer o.harvest(id)

	err := h.cmd.Wait()
	if err != nil {
		o.log.Error().Err(err).Msg("error waiting on job")
		// No return - continue.
	}

	endTime := time.Now()

	state := job.State{
		Status:       job.StatusDone,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	exitCode := h.cmd.ProcessState.ExitCode()
	state.ExitCode = &exitCode
	if *state.ExitCode != 0 {
		state.Status = job.StatusFailed
	}

	usage, err := process.GetUsage(h.cmd)
	if err != nil {
		o.log.Error().Err(err).Str("job", id).Msg("could not retrieve usage information")
	}
	state.ResourceUsage = usage

	return state, nil
}
