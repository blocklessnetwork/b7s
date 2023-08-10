package pbft

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) startViewChange(view uint) error {

	r.log.Info().Uint("current_view", r.view).Msg("starting view change")

	r.stopRequestTimer()

	r.view = view
	r.activeView = false

	vc := ViewChange{
		View:     r.view,
		Prepares: r.getPrepareSet(),
	}

	err := r.broadcast(vc)
	if err != nil {
		return fmt.Errorf("could not broadcast view change: %w", err)
	}

	r.log.Info().Uint("pending_view", r.view).Msg("view change successfully broadcast")

	return nil
}

func (r *Replica) processViewChange(replica peer.ID, msg ViewChange) error {

	log := r.log.With().Str("replica", replica.String()).Uint("view", msg.View).Logger()
	log.Info().Msg("processing view change message")

	if msg.View < r.view {
		log.Warn().Uint("current_view", r.view).Msg("received view change for an old view")
		return nil
	}

	// Check if the view change message is valid.
	err := r.validViewChange(msg)
	if err != nil {
		return fmt.Errorf("view change message is not valid (replica: %s): %w", replica.String(), err)
	}

	r.recordViewChangeReceipt(replica, msg)

	log.Info().Msg("processed view change message")

	nextView, should := r.shouldSendViewChange()
	if should {

		log.Info().Uint("next_view", nextView).Msg("we have received enough view change messages, joining view change")

		err = r.startViewChange(nextView)
		if err != nil {
			log.Error().Err(err).Uint("next_view", nextView).Msg("unable to send view change")
		}
	}

	projectedPrimary := r.peers[r.primary(msg.View)]
	log.Info().Str("id", projectedPrimary.String()).Msg("expected primary for the view")

	// If `I` am not the expected primary for this view - I've done all I should.
	if projectedPrimary != r.id {
		log.Info().Msg("processed view change message")
		return nil
	}

	// I am the primary for the view in question.
	if !r.viewChangeReady(msg.View) {
		log.Info().Msg("I am the expected primary for the view, but not enough view change messages yet")
		return nil
	}

	log.Info().Msg("I am the expected primary for the new view, have enough view change messages")

	return r.startNewView(msg.View)
}

func (r *Replica) recordViewChangeReceipt(replica peer.ID, vc ViewChange) {

	vcs, ok := r.viewChanges[vc.View]
	if !ok {
		r.viewChanges[vc.View] = newViewChangeReceipts()
		vcs = r.viewChanges[vc.View]
	}

	vcs.Lock()
	defer vcs.Unlock()

	_, exists := vcs.m[replica]
	if exists {
		r.log.Warn().Uint("view", vc.View).Str("replica", replica.String()).Msg("ignoring duplicate view change message")
		return
	}

	vcs.m[replica] = vc
}

// Required for a view change, getPrepareSet returns the set of all requests prepared on this replica.
// It includes a valid pre-prepare message and 2f matching, valid prepare messages signed by other backups - same view, sequence number and digest.
func (r *Replica) getPrepareSet() []PrepareInfo {

	r.log.Info().Msg("determining prepare set")

	r.log.Debug().Interface("state", r.replicaState).Msg("current state for the replica")

	var out []PrepareInfo

	for msgID, prepare := range r.prepares {

		log := r.log.With().Uint("view", msgID.view).Uint("sequence", msgID.sequence).Logger()

		for digest := range r.requests {

			log = log.With().Str("digest", digest).Logger()
			log.Info().Msg("checking if request is suitable for prepare set")

			if !r.prepared(msgID.view, msgID.sequence, digest) {
				log.Info().Msg("request not prepared - skipping")
				continue
			}

			log.Info().Msg("request prepared - including")

			prepareInfo := PrepareInfo{
				View:           msgID.view,
				SequenceNumber: msgID.sequence,
				Digest:         digest,
				PrePrepare:     r.preprepares[msgID],
				Prepares:       prepare.m,
			}

			out = append(out, prepareInfo)
		}
	}

	r.log.Debug().Interface("prepare_set", out).Msg("prepare set for the replica")

	return out
}

func (r *Replica) validViewChange(vc ViewChange) error {

	if vc.View == 0 {
		return errors.New("invalid view number")
	}

	for _, prepare := range vc.Prepares {

		if prepare.View >= vc.View || prepare.SequenceNumber == 0 {
			return fmt.Errorf("view change - prepare has an invalid view/sequence number (view: %v, prepare view: %v, sequence: %v)", vc.View, prepare.View, prepare.SequenceNumber)
		}

		if prepare.View != prepare.PrePrepare.View || prepare.SequenceNumber != prepare.PrePrepare.SequenceNumber {
			return fmt.Errorf("view change - prepare has an unmatching pre-prepare message (view/sequence number)")
		}

		if prepare.Digest == "" {
			return fmt.Errorf("view change - prepare has an empty digest")
		}

		if prepare.Digest != prepare.PrePrepare.Digest {
			return fmt.Errorf("view change - prepare has an unmatching pre-prepare message (digest)")
		}

		if uint(len(prepare.Prepares)) < r.prepareQuorum() {
			return fmt.Errorf("view change - prepare has an insufficient number of prepare messages (have: %v)", len(prepare.Prepares))
		}

		for _, pp := range prepare.Prepares {
			if pp.View != prepare.View || pp.SequenceNumber != prepare.SequenceNumber || pp.Digest != prepare.Digest {
				return fmt.Errorf("view change - included prepare message for wrong request")
			}
		}

	}

	return nil

	// TODO (pbft): Is this below relevant?

	// Condition:
	// !(p.View < vc.View && p.SequenceNumber > vc.H && p.SequenceNumber <= vc.H+instance.L)
	//
	// Translate to english:
	// view change is bad if condition is NOT met:
	// so, view is good if the following is true:
	// p.View < vc.View && p.SequenceNumber > vc.H && p.SequenceNumber <= vc.H+instance.L
	//
	// - prepares have to be for a view lower than one received
	// - prepares have a sequence number higher than the view change's H value
	// - prepares have a sequence number lower than the view change's H + L value
}

// Liveness condition - if we received f+1 valid view change messages from other replicas,
// (greater than our current view), send a view change message for the smallest view in the set. Do so
// even if our timer has not expired.
func (r *Replica) shouldSendViewChange() (uint, bool) {

	// If we're already participating in a view change, we're good.
	if !r.activeView {
		return 0, false
	}

	var newView uint

	// Go through view change messages we have received.
	for view, vcs := range r.viewChanges {
		// Only consider views higher than our current one.
		if view <= r.view {
			continue
		}

		vcs.Lock()

		// See how many view change messages we received. Don't count our own.
		count := 0
		for replica := range vcs.m {
			if replica != r.id {
				count++
			}

			// NOTE: We already check if the view change is valid on receiving it.
		}

		vcs.Unlock()

		// If we have more than f+1 view change messages, consider sending one too.
		if uint(count) >= r.f+1 {

			if newView == 0 { // Set new view if it was uninitialized.
				newView = view
			} else if view < newView { // We have multiple sets of f+1 view change messages. Use the lowest one.
				newView = view
			}
		}
	}

	if newView != 0 {
		return newView, true
	}

	return newView, false
}
