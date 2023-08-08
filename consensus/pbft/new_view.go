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

	preprepares := r.generatePreprepares(view, vcs.m)

	newView := NewView{
		View:        view,
		Messages:    vcs.m,
		PrePrepares: preprepares,
	}

	err := r.broadcast(newView)
	if err != nil {
		return fmt.Errorf("could not broadcast new-view message (view: %v): %w", view, err)
	}

	// Now, save any information we did not have previously (e.g. requests), change the current view for the replica and enter the view (set as active).
	for _, preprepare := range preprepares {
		_, found := r.requests[preprepare.Digest]
		if !found {
			r.requests[preprepare.Digest] = preprepare.Request
		}

		_, found = r.pending[preprepare.Digest]
		if !found {
			r.pending[preprepare.Digest] = preprepare.Request
		}
	}

	r.view = view
	r.activeView = true

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

func (r *Replica) processNewView(replica peer.ID, newView NewView) error {

	log := r.log.With().Str("replica", replica.String()).Uint("new_view", newView.View).Logger()

	log.Info().Msg("received new view message")

	if r.activeView {
		return ErrActiveView
	}

	if newView.View <= r.view {
		log.Warn().Uint("current_view", r.view).Msg("received new view message for a view lower than ours, discarding")
		return nil
	}

	// Make sure that the replica sending this is the replica that will be the primary for the view in question.
	projectedPrimary := r.peers[r.primary(newView.View)]
	if projectedPrimary != replica {
		return fmt.Errorf("projected primary for the view isn't the sender of the new-view message (projected: %v, sender: %v)",
			projectedPrimary.String(),
			replica.String())
	}

	// Verify number of messages included.
	count := uint(len(newView.Messages))
	haveQuorum := count >= r.commitQuorum()

	if !haveQuorum {
		return fmt.Errorf("new-view message does not have a quorum of view-change messages (replica: %v, count: %v)", replica.String(), count)
	}

	// Go through ViewChange messages and validate them.
	for _, vc := range newView.Messages {
		if vc.View != newView.View {
			return fmt.Errorf("view change message references a wrong view (view_change: %v, new_view: %v)", vc.View, newView.View)
		}

		// TODO (pbft): verify signatures.
	}

	for i, preprepare := range newView.PrePrepares {
		if preprepare.View != newView.View {
			return fmt.Errorf("new view preprepare message for a wrong view (preprepare_view: %v, new_view: %v)", preprepare.View, newView.View)
		}

		// Verify our sequence numbers are all there, though offset by one.
		if uint(i) != preprepare.SequenceNumber {
			log.Warn().Interface("preprepares", newView.PrePrepares).Msg("preprepares have unexpected sequence number value (possible gap)")
			return fmt.Errorf("unexpected sequence number gap")
		}
	}

	// Update our local view, switch to active view.
	r.view = newView.View
	r.activeView = true

	// Start processing preprepares.
	for _, preprepare := range newView.PrePrepares {
		err := r.processPrePrepare(replica, preprepare)
		if err != nil {
			log.Error().Err(err).Uint("view", preprepare.View).Uint("sequence", preprepare.SequenceNumber).Msg("error processing preprepare message")
			// Continue despite errors.
		}
	}

	log.Info().Msg("processed new view message")

	return nil
}
