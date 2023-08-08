package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) sendPrePrepare(req Request) error {

	// Only primary replica can send pre-prepares.
	if !r.isPrimary() {
		return nil
	}

	r.startRequestTimer(false)

	seqNo := r.sequence + 1
	r.sequence++

	msg := PrePrepare{
		View:           r.view,
		SequenceNumber: seqNo,
		Request:        req,
		Digest:         getDigest(req),
	}

	log := r.log.With().Uint("view", msg.View).Uint("sequence_number", msg.SequenceNumber).Str("digest", msg.Digest).Logger()

	// TODO (pbft): Check if we had this or other pre-prepares for this request.
	if r.conflictingPrePrepare(msg) {
		return fmt.Errorf("dropping pre-prepare as we have a conflicting one")
	}

	log.Info().Msg("broadcasting pre-prepare message")

	err := r.broadcast(msg)
	if err != nil {
		return fmt.Errorf("could not broadcast pre-prepare message: %w", err)
	}

	log.Info().Msg("pre-prepare message successfully broadcast")

	// Take a note of this pre-prepare. This will naturally only happen on the primary replica.
	r.preprepares[getMessageID(msg.View, msg.SequenceNumber)] = msg

	return nil
}

// Process a pre-prepare message. This should naturally only happen on non-primary replicas.
func (r *Replica) processPrePrepare(replica peer.ID, msg PrePrepare) error {

	if r.isPrimary() {
		r.log.Warn().Msg("primary replica received a pre-prepare, dropping")
		return nil
	}

	log := r.log.With().Str("replica", replica.String()).Uint("view", msg.View).Uint("sequence_no", msg.SequenceNumber).Str("digest", msg.Digest).Logger()

	log.Info().Msg("received pre-prepare message")

	if !r.activeView {
		return ErrViewChange
	}

	if replica != r.primaryReplicaID() {
		log.Warn().Str("primary", r.primaryReplicaID().String()).Msg("pre-prepare came from a replica that is not the primary, dropping")
		return nil
	}

	if msg.View != r.view {
		return fmt.Errorf("pre-prepare has an invalid view value (received: %v, current: %v)", msg.View, r.view)
	}

	id := getMessageID(msg.View, msg.SequenceNumber)

	// TODO (pbft): in reality more involved, for now we'll stop if there's something existing already.
	existing, ok := r.preprepares[id]
	if ok {
		log.Warn().Str("existing_digest", existing.Digest).Msg("pre-prepare message already exists for this view and sequence number, dropping")
		return nil
	}

	// We don't have this pre-prepare. Save it now.
	r.preprepares[id] = msg

	// Save this request.
	r.requests[msg.Digest] = msg.Request
	r.pending[msg.Digest] = msg.Request

	if !r.prePrepared(msg.View, msg.SequenceNumber, msg.Digest) {
		log.Warn().Msg("request is not pre-prepared, stopping")
		return nil
	}

	// Broadcast prepare message.
	err := r.sendPrepare(msg)
	if err != nil {
		return fmt.Errorf("could not send prepare message: %w", err)
	}

	return nil
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
