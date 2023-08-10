package pbft

import (
	"time"
)

func (r *Replica) startRequestTimer(overrideExisting bool) {

	r.log.Info().Msg("starting view change timer")

	if r.requestTimer != nil && !overrideExisting {
		r.log.Info().Msg("view change timer running, not overriding")
		return
	}

	if r.requestTimer != nil {
		r.requestTimer.Stop()
	}

	r.requestTimer = time.AfterFunc(RequestTimeout, func() {
		r.sl.Lock()
		defer r.sl.Unlock()

		err := r.startViewChange(r.view + 1)
		if err != nil {
			r.log.Error().Err(err).Msg("could not start view change")
		}
	})

	r.log.Info().Msg("view change timer started")
}

func (r *Replica) stopRequestTimer() {

	r.log.Info().Msg("stopping view change timer")

	if r.requestTimer == nil {
		r.log.Info().Msg("no active view change timer")
		return
	}

	r.requestTimer.Stop()
	r.requestTimer = nil

	r.log.Info().Msg("view change timer stopped")
}
