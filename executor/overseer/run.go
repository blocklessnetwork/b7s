package overseer

import (
	"fmt"
)

func (o *Overseer) Run(job Job) (JobState, error) {

	h, err := o.Start(job)
	if err != nil {
		return JobState{}, fmt.Errorf("could not start job: %w", err)
	}

	state, err := o.Wait(h.ID)
	if err != nil {
		return JobState{}, fmt.Errorf("could not wait on job: %w", err)
	}

	return state, nil
}
