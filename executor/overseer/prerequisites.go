package overseer

import (
	"fmt"
)

func (o *Overseer) checkPrerequisites(job Job) error {

	// If we have an allowlist, check if the job fits the criteria.
	if len(o.cfg.Allowlist) > 0 {
		allowed := false
		for _, exe := range o.cfg.Allowlist {
			if exe == job.Exec.Path {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("job executable is not in the allowlist (exe: %v)", job.Exec.Path)
		}
	}

	return nil
}
