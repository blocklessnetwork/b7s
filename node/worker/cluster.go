package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/request"
)

func (w *Worker) processFormCluster(ctx context.Context, from peer.ID, req request.FormCluster) error {

	w.Log().Info().
		Str("request", req.RequestID).
		Strs("peers", blockless.PeerIDsToStr(req.Peers)).
		Stringer("consensus", req.Consensus).
		Msg("received request to form consensus cluster")

	// Add connection info about peers if we're not already connected to them.
	for _, addrInfo := range req.ConnectionInfo {

		if w.Host().ID() == addrInfo.ID {
			continue
		}

		if w.Connected(addrInfo.ID) {
			continue
		}

		w.Log().Debug().
			Any("known", w.Host().Network().Peerstore().Addrs(addrInfo.ID)).
			Any("received", addrInfo.Addrs).
			Stringer("peer", addrInfo.ID).
			Msg("received addresses for fellow cluster replica")

		w.Host().Network().Peerstore().AddAddrs(addrInfo.ID, addrInfo.Addrs, ClusterAddressTTL)
	}

	switch req.Consensus {
	case consensus.Raft:
		return w.createRaftCluster(ctx, from, req)

	case consensus.PBFT:
		return w.createPBFTCluster(ctx, from, req)
	}

	return fmt.Errorf("invalid consensus specified (%v %s)", req.Consensus, req.Consensus.String())
}

// processDisbandCluster will start cluster shutdown command.
func (w *Worker) processDisbandCluster(ctx context.Context, from peer.ID, req request.DisbandCluster) error {

	w.Log().Info().
		Stringer("peer", from).
		Str("request", req.RequestID).
		Msg("received request to disband consensus cluster")

	err := w.leaveCluster(req.RequestID, consensusClusterDisbandTimeout)
	if err != nil {
		return fmt.Errorf("could not disband cluster (request: %s): %w", req.RequestID, err)
	}

	w.Log().Info().
		Stringer("peer", from).
		Str("request", req.RequestID).
		Msg("left consensus cluster")

	return nil
}

func (w *Worker) leaveCluster(requestID string, timeout time.Duration) error {

	// Shutdown can take a while so use short locking intervals.
	cluster, ok := w.clusters.Get(requestID)
	if !ok {
		return errors.New("no cluster with that ID")
	}

	// TODO: Fix this logging.
	w.Log().Info().
		Stringer("consensus", cluster.Consensus()).
		Str("request", requestID).
		Msg("leaving consensus cluster")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// We know that the request is done executing when we have a result for it.
	_, ok = w.executeResponses.WaitFor(ctx, requestID)

	log := w.Log().With().Str("request", requestID).Logger()
	log.Info().Bool("executed_work", ok).Msg("waiting for execution done, leaving cluster")

	err := cluster.Shutdown()
	if err != nil {
		// Not much we can do at this point.
		return fmt.Errorf("could not leave cluster (request: %v): %w", requestID, err)
	}

	w.clusters.Delete(requestID)

	return nil
}
