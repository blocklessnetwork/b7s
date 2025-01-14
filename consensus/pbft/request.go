package pbft

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/blockless"
)

func (r *Replica) processRequest(ctx context.Context, from peer.ID, req Request) error {

	pub, err := req.Origin.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf("could not extract public key from client: %w", err)
	}

	err = req.Execute.VerifySignature(pub)
	if err != nil {
		return fmt.Errorf("request is not properly signed by the client: %w", err)
	}

	digest := getDigest(req)

	log := r.log.With().Str("client", from.String()).Str("request", req.ID).Str("digest", digest).Logger()

	log.Info().Msg("received a request")

	// Check if we've executed this before. If yes, just return the result.
	result, ok := r.executions[req.ID]
	if ok {
		log.Info().Msg("request already executed, sending result to client")

		err := r.send(ctx, req.Origin, &result, blockless.ProtocolID)
		if err != nil {
			return fmt.Errorf("could not send execution result back to client (request: %s, client: %s): %w", req.ID, req.Origin.String(), err)
		}

		return nil
	}

	// If we're not the primary, we'll drop the request. We do start a request timer though.
	if !r.isPrimary() {
		r.startRequestTimer(false)
		log.Info().Str("primary", r.primaryReplicaID().String()).Msg("we are not the primary replica, dropping the request")

		// Just to be safe, store the request we've seen.
		r.requests[digest] = req
		return nil
	}

	log.Info().Msg("we are the primary, processing the request")

	_, pending := r.pending[digest]
	if pending {
		return fmt.Errorf("this request is already queued, dropping (request: %v)", req.ID)
	}

	_, seen := r.requests[digest]
	if seen {
		log.Info().Msg("already seen this request, resubmitted")
	}

	// Take a note of this request.
	r.requests[digest] = req
	r.pending[digest] = req

	// Broadcast a pre-prepare message.
	err = r.sendPrePrepare(ctx, req)
	if err != nil {
		return fmt.Errorf("could not broadcast pre-prepare message (request: %v): %w", req.ID, err)
	}

	log.Info().Msg("processed request")

	return nil
}
