package pbft

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) startNewView(ctx context.Context, view uint) error {

	log := r.log.With().Uint("view", view).Logger()

	log.Info().Msg("starting a new view")

	projectedPrimary := r.peers[r.primary(view)]
	if projectedPrimary != r.id {
		return fmt.Errorf("am not the expected primary for the specified view (view: %v, primary: %v)", view, projectedPrimary.String())
	}

	vcs, ok := r.viewChanges[view]
	if !ok {
		return fmt.Errorf("no view change messages for the specified view (view: %v)", view)
	}

	// If we don't have our own view change message added yet - do it now.
	// Don't defer unlock because we invoke viewChangeReady, which locks the same view change slot.
	vcs.Lock()
	_, ok = vcs.m[r.id]
	if !ok {

		vc := ViewChange{
			View:     view,
			Prepares: r.getPrepareSet(),
		}

		vcs.m[r.id] = vc
	}
	vcs.Unlock()

	// Recheck that we have a valid view change state (quorum).
	if !r.viewChangeReady(view) {
		return fmt.Errorf("new view sequence started but not enough view change messages present (view: %v)", view)
	}

	log.Info().Msg("view change ready, broadcasting new view message")

	preprepares := r.generatePreprepares(view, vcs.m)

	for i := 0; i < len(preprepares); i++ {
		err := r.sign(&preprepares[i])
		if err != nil {
			return fmt.Errorf("new-view - could not sign preprepare message: %w", err)
		}
	}

	newView := NewView{
		View:        view,
		Messages:    vcs.m,
		PrePrepares: preprepares,
	}

	err := r.sign(&newView)
	if err != nil {
		return fmt.Errorf("could not sign the new view message: %w", err)
	}

	err = r.broadcast(newView)
	if err != nil {
		return fmt.Errorf("could not broadcast new-view message (view: %v): %w", view, err)
	}

	log.Info().Interface("new_view", newView).Msg("new view message successfully broadcast")

	r.cleanupState(view)

	// Now, save any information we did not have previously (e.g. new pre-prepares, requests), change the current view for the replica and enter the view (set as active).
	for _, preprepare := range preprepares {

		r.preprepares[getMessageID(preprepare.View, preprepare.SequenceNumber)] = preprepare

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

	log.Info().Msg("new view started")

	// See any pending requests you've seen and add them to the pipeline.
	outstandingRequests := r.outstandingRequests()
	for _, request := range outstandingRequests {

		_, pending := r.pending[getDigest(request)]
		if pending {
			r.log.Debug().Str("request", request.ID).Msg("outstanding request but already pending as part of a new-view payload, skipping")
			continue
		}

		err = r.processRequest(ctx, request.Origin, request)
		if err != nil {
			r.log.Error().Err(err).Str("request", request.ID).Msg("could not process request")
			// Log but continue.
		}
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

	// Phase 2. Go through all sequence numbers from 1 to max. If there is a prepare certificate
	// for a sequence number in the view change messages - create a pre-prepare message for m,v+1,n.
	// If there are multiple prepare certificates with different view numbers - use the highest view number.
	preprepares := make([]PrePrepare, 0, max)
	for sequenceNo := uint(1); sequenceNo <= max; sequenceNo++ {

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
			if !found || (prepare.View > out.View) {
				out = prepare
				// Could also compare with an empty PrepareInfo tbh.
				found = true
			}
		}
	}

	return out, found
}

func (r *Replica) processNewView(ctx context.Context, replica peer.ID, newView NewView) error {

	log := r.log.With().Str("replica", replica.String()).Uint("new_view", newView.View).Logger()

	log.Info().Interface("new_view", newView).Msg("processing new view message")

	if newView.View < r.view {
		log.Warn().Uint("current_view", r.view).Msg("received new view message for a view lower than ours, discarding")
		return nil
	}

	// Make sure that the replica sending this is the replica that will be the primary for the view in question.
	projectedPrimary := r.peers[r.primary(newView.View)]
	if projectedPrimary != replica {
		return fmt.Errorf("sender of the new-view message is not the projected primary for the view (sender: %v, projected: %v)",
			replica.String(),
			projectedPrimary.String())
	}

	err := r.verifySignature(&newView, projectedPrimary)
	if err != nil {
		return fmt.Errorf("could not verify signature of the new view message: %w", err)
	}

	// Verify number of messages included.
	count := uint(len(newView.Messages))
	haveQuorum := count >= r.commitQuorum()

	if !haveQuorum {
		return fmt.Errorf("new-view message does not have a quorum of view-change messages (sender: %v, count: %v)", replica.String(), count)
	}

	for replica, vc := range newView.Messages {

		err := r.validViewChange(vc)
		if err != nil {
			return fmt.Errorf("new-view - included view change message is invalid: %w", err)
		}

		err = r.verifySignature(&vc, replica)
		if err != nil {
			return fmt.Errorf("new-view - included view change message has an invalid signature: %w", err)
		}
	}

	for i, preprepare := range newView.PrePrepares {
		if preprepare.View != newView.View {
			return fmt.Errorf("new view preprepare message for a wrong view (preprepare_view: %v, new_view: %v)", preprepare.View, newView.View)
		}

		// Verify sequence numbers are all there, though offset by one.
		if uint(i+1) != preprepare.SequenceNumber {
			log.Warn().Interface("preprepares", newView.PrePrepares).Msg("preprepares have unexpected sequence number value (possible gap)")
			return fmt.Errorf("unexpected sequence number list")
		}
	}

	// Update our local view, switch to active view.
	r.view = newView.View
	r.activeView = true

	log.Info().Msg("processed new view message")

	r.log.Info().Str("primary", r.primaryReplicaID().String()).Uint("view", r.view).Msg("entered new view")

	r.cleanupState(newView.View)

	// Start processing preprepares.
	for _, preprepare := range newView.PrePrepares {
		err := r.processPrePrepare(ctx, replica, preprepare)
		if err != nil {
			log.Error().Err(err).Uint("view", preprepare.View).Uint("sequence", preprepare.SequenceNumber).Msg("error processing preprepare message")
			// Continue despite errors.
		}
	}

	log.Info().Msg("processed preprepares from the new view")

	if len(r.outstandingRequests()) > 0 {
		r.log.Info().Msg("outstanding requests found, starting a new view change timer")
		r.startRequestTimer(false)
	}

	return nil
}
