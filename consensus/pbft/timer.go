package pbft

import (
	"time"
)

func (r *Replica) startRequestTimer(overrideExisting bool) {

	r.log.Info().Msg("starting view change timer")

	if r.requestTimer != nil && !overrideExisting {
		r.log.Info().Bool("override", overrideExisting).Msg("view change timer running, not overriding")
		return
	}

	// TODO (pbft): Proper stopping/draining of the timer.
	r.requestTimer.Stop()
	r.requestTimer = time.NewTimer(RequestTimeout)

	go func() {
		<-r.requestTimer.C

		r.sl.Lock()
		defer r.sl.Unlock()

		err := r.startViewChange(r.view + 1)
		if err != nil {
			r.log.Error().Err(err).Msg("could not start view change")
		}
	}()
}

func (r *Replica) stopRequestTimer() {

	r.log.Info().Msg("stopping view change timer")

	if r.requestTimer == nil {
		r.log.Info().Msg("no active view change timmer")
		return
	}

	r.requestTimer.Stop()
	r.requestTimer = nil
}
