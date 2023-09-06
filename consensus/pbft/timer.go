package pbft

import (
	"time"
)

func (r *Replica) startRequestTimer(overrideExisting bool) {

	r.log.Debug().Msg("starting view change timer")

	if r.requestTimer != nil && !overrideExisting {
		r.log.Debug().Msg("view change timer running, not overriding")
		return
	}

	if r.requestTimer != nil {
		r.requestTimer.Stop()
	}

	// Evaluate the view number now. Potentially, we could've already advanced
	// to the next view before our inactivity timer fires.
	targetView := r.view + 1

	r.requestTimer = time.AfterFunc(r.cfg.RequestTimeout, func() {
		r.sl.Lock()
		defer r.sl.Unlock()

		err := r.startViewChange(targetView)
		if err != nil {
			r.log.Error().Err(err).Msg("could not start view change")
		}
	})

	r.log.Debug().Msg("view change timer started")
}

func (r *Replica) stopRequestTimer() {

	r.log.Debug().Msg("stopping view change timer")

	if r.requestTimer == nil {
		r.log.Debug().Msg("no active view change timer")
		return
	}

	r.requestTimer.Stop()
	r.requestTimer = nil

	r.log.Debug().Msg("view change timer stopped")
}
