package overseer

import (
	"fmt"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

func (o *Overseer) Run(j job.Job) (job.State, error) {

	id, err := o.Start(j)
	if err != nil {
		return job.State{}, fmt.Errorf("could not start job: %w", err)
	}

	state, err := o.Wait(id)
	if err != nil {
		return job.State{}, fmt.Errorf("could not wait on job: %w", err)
	}

	return state, nil
}
