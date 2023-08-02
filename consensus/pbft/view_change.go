package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) startViewChange() error {

	r.sl.Lock()
	defer r.sl.Unlock()

	r.log.Info().Uint("current_view", r.view).Msg("starting view change")

	r.stopRequestTimer()

	r.view++
	r.activeView = false

	// TODO (pbft): Signing.
	// TODO (pbft): Resending.

	vc := ViewChange{
		View:     r.view,
		Prepares: r.getPrepareSet(),
	}

	err := r.broadcast(vc)
	if err != nil {
		return fmt.Errorf("could not broadcast view change: %w", err)
	}

	return nil
}

func (r *Replica) processViewChange(replica peer.ID, msg ViewChange) error {

	r.log.Info().Str("replica", replica.String()).Msg("processing view change message")

	if msg.View < r.view {
		r.log.Warn().Uint("current", r.view).Uint("received", msg.View).Msg("received view change for an old view")
		return nil
	}

	return fmt.Errorf("TBD: not implemented")
}

// Required for a view change, getPrepareSet returns the set of all requests prepared on this replica.
// It includes a valid pre-prepare message and 2f matching, valid prepare messages signed by other backups - same view, sequence number and digest.
func (r *Replica) getPrepareSet() []PrepareInfo {

	r.log.Info().Msg("determining prepare set")

	r.log.Debug().Interface("state", r.replicaState).Msg("current state for the replica")

	var out []PrepareInfo

	for msgID, prepare := range r.prepares {

		log := r.log.With().Uint("view", msgID.view).Uint("sequnce", msgID.sequence).Logger()

		for digest := range r.requests {

			log = log.With().Str("digest", digest).Logger()
			log.Info().Msg("checking if request is suitable for prepare set")

			if !r.prepared(msgID.view, msgID.sequence, digest) {
				log.Info().Msg("request not prepared - skipping")
				continue
			}

			log.Info().Msg("request prepared - including")

			prepareInfo := PrepareInfo{
				View:       msgID.view,
				Sequnce:    msgID.sequence,
				Digest:     digest,
				PrePrepare: r.preprepares[msgID],
				Prepares:   prepare.m,
			}

			out = append(out, prepareInfo)
		}
	}

	r.log.Debug().Interface("prepare_set", out).Msg("prepare set for the replica")

	return out
}
