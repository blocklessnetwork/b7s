package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) startNewView(view uint) error {

	vcs, ok := r.viewChanges[view]
	if !ok {
		return fmt.Errorf("no view change messages for the specified view (view: %v)", view)
	}

	vcs.Lock()
	defer vcs.Unlock()

	// TODO (pbft): If we don't have our own view change message yet - add it here now.

	// Recheck that we have a valid view change state (quorum).
	if !r.viewChangeReady(view) {
		return fmt.Errorf("new view sequence started but not enough view change messages present (view: %v)", view)
	}

	msgs := make([]ViewChange, 0, len(vcs.m))
	for _, vc := range vcs.m {
		msgs = append(msgs, vc)
	}

	preprepares := r.generatePreprepares(view, vcs.m)

	newView := NewView{
		View:        view,
		Messages:    msgs,
		PrePrepares: preprepares,
	}

	err := r.broadcast(newView)
	if err != nil {
		return fmt.Errorf("could not broadcast new-view message (view: %v): %w", view, err)
	}

	return nil
}

func (r *Replica) generatePreprepares(view uint, vcs map[peer.ID]ViewChange) []PrePrepare {

	log := r.log.With().Uint("view", view).Logger()

	// Phase 1. We don't have checkpoints, so our lower sequence number bound is 0.
	// Determine the upper higher sequence bound by going through the view change messages
	// and examining the prepare certificates.
	max := getHighestSequenceNumber(vcs)

	log.Info().Uint("max", max).Msg("generating preprepares for new view, determined max sequence number")

	// Phase 2. Go through all sequence numbers from 0 to max. If there is a prepare certificate
	// for a sequence number in the view change messages - create a pre-prepare message for m,v+1,n.
	// If there are multiple prepare certificates with different view numbers - use the highest view number.
	preprepares := make([]PrePrepare, 0, max)
	for sequenceNo := uint(0); sequenceNo <= max; sequenceNo++ {

		log := log.With().Uint("sequence", sequenceNo).Logger()

		prepare, exists := getPrepare(vcs, sequenceNo)
		// If we have a prepare certificate for this sequence number, add it.
		if exists {

			log.Info().Str("digest", prepare.Digest).Str("request", prepare.PrePrepare.Request.ID).Msg("generating preprepares for new view, found prepare certificate")

			preprepare := PrePrepare{
				View:           view,
				SequenceNumber: sequenceNo,
				Digest:         prepare.Digest,
				Request:        prepare.PrePrepare.Request,
			}

			preprepares = append(preprepares, preprepare)
			continue
		}

		log.Info().Msg("generating preprepares for new view, no prepare certificate found, using a null request")

		// We don't have a prepare certificate for this sequence number - create a preprepare for a null request.
		preprepare := PrePrepare{
			View:           view,
			SequenceNumber: sequenceNo,
			Digest:         "",
			Request:        NullRequest,
		}

		preprepares = append(preprepares, preprepare)
	}

	return preprepares
}

func getHighestSequenceNumber(vcs map[peer.ID]ViewChange) uint {

	var max uint

	// For each view change message (from a replica).
	for _, vc := range vcs {
		// Go through all prepares.
		for _, prepare := range vc.Prepares {
			// Update the max sequence number seen if current one is higher.
			if prepare.SequenceNumber > max {
				max = prepare.SequenceNumber
			}
		}
	}

	return max
}

func getPrepare(vcs map[peer.ID]ViewChange, sequenceNo uint) (PrepareInfo, bool) {

	var (
		out   PrepareInfo
		found bool
	)

	// For each view change message (from a replica).
	for _, vc := range vcs {

		// Go through prepares.
		for _, prepare := range vc.Prepares {

			// Only observe the prepares with this sequence number.
			if prepare.SequenceNumber != sequenceNo {
				continue
			}

			// In case we have multiple prepares for the same sequence number,
			// keep the one from the highest view.
			if prepare.View > out.View {
				out = prepare
				// Could also compare with an empty PrepareInfo tbh.
				found = true
			}
		}
	}

	return out, found
}
