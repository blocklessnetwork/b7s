package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) processRequest(from peer.ID, req Request) error {

	r.log.Info().Str("peer", from.String()).Str("id", req.ID).Msg("received a request")

	if !r.activeView {
		return ErrViewChange
	}

	// If we're not the primary, we'll drop the request. We do start a request timer though.
	if !r.isPrimary() {
		r.startRequestTimer(false)
		r.log.Warn().Str("primary", r.primaryReplicaID().String()).Msg("we are not the primary replica, dropping the request")
		return nil
	}

	digest := getDigest(req)

	r.log.Info().Str("id", req.ID).Str("digest", digest).Msg("we are the primary, processing the request")

	_, found := r.requests[digest]
	if found {
		return fmt.Errorf("already seen this request, dropping")
	}

	// Take a note of this request.
	r.requests[digest] = req
	r.pending[digest] = req

	// Broadcast a pre-prepare message.
	err := r.sendPrePrepare(req)
	if err != nil {
		return fmt.Errorf("could not broadcast pre-prepare message: %w", err)
	}

	return nil
}
