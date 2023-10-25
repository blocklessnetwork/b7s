package overseer

import (
	"fmt"
)

// checkPrerequisites will check if the job is allowed to run - e.g. if the allowlist or denylist exclude it.
func (o *Overseer) checkPrerequisites(job Job) error {

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

	if len(o.cfg.Denylist) > 0 {
		denied := false
		for _, exe := range o.cfg.Denylist {
			if exe == job.Exec.Path {
				denied = true
				break
			}
		}

		if denied {
			return fmt.Errorf("job executable is in the denylist (exe: %v)", job.Exec.Path)
		}
	}

	return nil
}
