package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) processRequest(from peer.ID, req Request) error {

	digest := getDigest(req)

	log := r.log.With().Str("client", from.String()).Str("request", req.ID).Str("digest", digest).Logger()

	log.Info().Msg("received a request")

	// If we're not the primary, we'll drop the request. We do start a request timer though.
	if !r.isPrimary() {
		r.startRequestTimer(false)
		log.Warn().Str("primary", r.primaryReplicaID().String()).Msg("we are not the primary replica, dropping the request")

		// Just to be safe, store the request we've seen.
		r.requests[digest] = req
		return nil
	}

	log.Info().Msg("we are the primary, processing the request")

	_, found := r.requests[digest]
	if found {
		return fmt.Errorf("already seen this request, dropping (request: %v)", req.ID)
	}

	// Take a note of this request.
	r.requests[digest] = req
	r.pending[digest] = req

	// Broadcast a pre-prepare message.
	err := r.sendPrePrepare(req)
	if err != nil {
		return fmt.Errorf("could not broadcast pre-prepare message (request: %v): %w", req.ID, err)
	}

	log.Info().Msg("processed request")

	return nil
}
