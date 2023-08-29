package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// Send a prepare message. Naturally this is only sent by the non-primary replicas,
// as a response to a pre-prepare message.
func (r *Replica) sendPrepare(preprepare PrePrepare) error {

	msg := Prepare{
		View:           preprepare.View,
		SequenceNumber: preprepare.SequenceNumber,
		Digest:         preprepare.Digest,
	}

	log := r.log.With().Str("digest", msg.Digest).Uint("view", msg.View).Uint("sequence_number", msg.SequenceNumber).Logger()

	log.Info().Msg("broadcasting prepare message")

	err := r.sign(&msg)
	if err != nil {
		return fmt.Errorf("could not sign prepare message: %w", err)
	}

	err = r.broadcast(msg)
	if err != nil {
		return fmt.Errorf("could not broadcast prepare message: %w", err)
	}

	log.Info().Msg("prepare message successfully broadcast")

	// Record this prepare message.
	r.recordPrepareReceipt(r.id, msg)

	return nil
}

func (r *Replica) recordPrepareReceipt(replica peer.ID, prepare Prepare) {

	msgID := getMessageID(prepare.View, prepare.SequenceNumber)
	prepares, ok := r.prepares[msgID]
	if !ok {
		r.prepares[msgID] = newPrepareReceipts()
		prepares = r.prepares[msgID]
	}

	prepares.Lock()
	defer prepares.Unlock()

	_, exists := prepares.m[replica]
	if exists {
		r.log.Warn().Uint("view", prepare.View).Uint("sequence", prepare.SequenceNumber).Str("digest", prepare.Digest).Str("replica", replica.String()).Msg("ignoring duplicate prepare message")
		return
	}

	prepares.m[replica] = prepare
}

func (r *Replica) processPrepare(replica peer.ID, prepare Prepare) error {

	log := r.log.With().Str("replica", replica.String()).Uint("view", prepare.View).Uint("sequence_no", prepare.SequenceNumber).Str("digest", prepare.Digest).Logger()

	log.Info().Msg("received prepare message")

	if replica == r.primaryReplicaID() {
		log.Warn().Msg("received prepare message from primary, ignoring")
		return nil
	}

	if prepare.View != r.view {
		return fmt.Errorf("prepare has an invalid view value (received: %v, current: %v)", prepare.View, r.view)
	}

	err := r.verifySignature(&prepare, replica)
	if err != nil {
		return fmt.Errorf("could not verify signature for the prepare message: %w", err)
	}

	r.recordPrepareReceipt(replica, prepare)

	log.Info().Msg("processed prepare message")

	return r.maybeSendCommit(prepare.View, prepare.SequenceNumber, prepare.Digest)
}
