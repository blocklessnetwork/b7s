package overseer

import (
	"errors"
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
	"github.com/blocklessnetwork/b7s/execution/overseer/overseer/internal/process"
)

func (o *Overseer) Kill(id string) (job.State, error) {

	o.Lock()
	h, ok := o.jobs[id]
	o.Unlock()

	if !ok {
		return job.State{}, errors.New("unknown job")
	}

	h.Lock()
	defer h.Unlock()

	defer o.harvest(id)

	if h.cmd.Process == nil {
		return job.State{}, errors.New("job is not running")
	}

	// TODO: Check if the process has started any processes.
	err := h.cmd.Process.Kill()
	if err != nil {
		return job.State{}, fmt.Errorf("could not kill process: %w", err)
	}

	endTime := time.Now()

	// Until we wait on the process (read its state) it will still exist on the system as a zombie process.
	err = h.cmd.Wait()
	if err != nil {

		signaled, _ := wasSignalled(h.cmd.ProcessState)
		if !signaled {
			return job.State{}, fmt.Errorf("could not wait on process: %w", err)
		}

		o.log.Trace().Err(err).Str("job", id).Msg("expected error - job was signaled, wait produced an error")
	}

	state := job.State{
		Status:       job.StatusKilled,
		Stdout:       h.stdout.String(),
		Stderr:       h.stderr.String(),
		StartTime:    h.start,
		EndTime:      &endTime,
		ObservedTime: time.Now(),
	}

	exitCode, status := determineProcessStatus(h.cmd.ProcessState)
	state.ExitCode = exitCode
	state.Status = status

	usage, err := process.GetUsage(h.cmd)
	if err != nil {
		o.log.Error().Err(err).Str("job", id).Msg("could not retrieve usage information")
	}
	state.ResourceUsage = usage

	return state, nil
}
