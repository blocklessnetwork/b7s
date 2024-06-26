package pbft

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) sendPrePrepare(req Request) error {

	// Only primary replica can send pre-prepares.
	if !r.isPrimary() {
		return nil
	}

	r.sequence++
	sequence := r.sequence

	msg := PrePrepare{
		View:           r.view,
		SequenceNumber: sequence,
		Request:        req,
		Digest:         getDigest(req),
	}

	log := r.log.With().Uint("view", msg.View).Uint("sequence_number", msg.SequenceNumber).Str("digest", msg.Digest).Logger()

	if r.conflictingPrePrepare(msg) {
		return fmt.Errorf("dropping pre-prepare as we have a conflicting one")
	}

	err := r.sign(&msg)
	if err != nil {
		return fmt.Errorf("could not sign pre-prepare message: %w", err)
	}

	log.Info().Msg("broadcasting pre-prepare message")

	err = r.broadcast(msg)
	if err != nil {
		return fmt.Errorf("could not broadcast pre-prepare message: %w", err)
	}

	log.Info().Msg("pre-prepare message successfully broadcast")

	// Take a note of this pre-prepare. This will naturally only happen on the primary replica.
	r.preprepares[getMessageID(msg.View, msg.SequenceNumber)] = msg

	return nil
}

// Process a pre-prepare message. This should only happen on backup replicas.
func (r *Replica) processPrePrepare(ctx context.Context, replica peer.ID, msg PrePrepare) error {

	if r.isPrimary() {
		r.log.Warn().Msg("primary replica received a pre-prepare, dropping")
		return nil
	}

	log := r.log.With().Str("replica", replica.String()).Uint("view", msg.View).Uint("sequence_no", msg.SequenceNumber).Str("digest", msg.Digest).Logger()

	log.Info().Msg("received pre-prepare message")

	if replica != r.primaryReplicaID() {
		log.Error().Str("primary", r.primaryReplicaID().String()).Msg("pre-prepare came from a replica that is not the primary, dropping")
		return nil
	}

	if msg.View != r.view {
		return fmt.Errorf("pre-prepare for an invalid view (received: %v, current: %v)", msg.View, r.view)
	}

	err := r.verifySignature(&msg, r.primaryReplicaID())
	if err != nil {
		return fmt.Errorf("pre-prepare message signature not valid: %w", err)
	}

	id := getMessageID(msg.View, msg.SequenceNumber)

	existing, ok := r.preprepares[id]
	if ok {
		log.Error().Str("existing_digest", existing.Digest).Msg("pre-prepare message already exists for this view and sequence number, dropping")
		return ErrConflictingPreprepare
	}

	// We don't have this pre-prepare. Save it now.
	r.preprepares[id] = msg

	// TODO (pbft): See if this is the same request we saw. If it isn't consider triggering a view change right here and now.
	// Save this request.
	r.requests[msg.Digest] = msg.Request
	r.pending[msg.Digest] = msg.Request

	r.startRequestTimer(false)

	// Just a sanity check at this point, since we've set up the state just now.
	if !r.prePrepared(msg.View, msg.SequenceNumber, msg.Digest) {
		log.Warn().Msg("request is not pre-prepared, stopping")
		return nil
	}

	log.Info().Msg("processed pre-prepare")

	// Broadcast prepare message.
	err = r.sendPrepare(msg)
	if err != nil {
		return fmt.Errorf("could not send prepare message: %w", err)
	}

	// There's a possibility our prepare was the one that pushes us into the quorum
	// and we now have the commit condition achieved.
	return r.maybeSendCommit(msg.View, msg.SequenceNumber, msg.Digest)
}

func (r *Replica) conflictingPrePrepare(preprepare PrePrepare) bool {

	for _, pp := range r.preprepares {

		// If we have a pre-prepare with the same view and same digest but different sequence number - invalid.
		if pp.View == preprepare.View && pp.Digest == preprepare.Digest && pp.SequenceNumber != preprepare.SequenceNumber {
			return true
		}
	}

	return false
}
