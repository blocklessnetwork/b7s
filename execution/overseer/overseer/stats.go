package overseer

import (
	"errors"
	"time"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
	"github.com/blocklessnetwork/b7s/execution/overseer/overseer/internal/process"
)

func (o *Overseer) Stats(id string) (job.State, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return job.State{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	state := job.State{
		Status:       job.StatusRunning,
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

	usage, err := process.GetUsage(h.cmd)
	if err != nil {
		o.log.Error().Err(err).Str("job", id).Msg("could not retrieve usage information")
	}
	state.ResourceUsage = usage

	return state, nil
}
