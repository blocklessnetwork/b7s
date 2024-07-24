package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog/log"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processFormCluster(ctx context.Context, from peer.ID, req request.FormCluster) error {

	// Should never happen.
	if !n.isWorker() {
		n.log.Warn().Str("peer", from.String()).Msg("only worker nodes participate in consensus clusters")
		return nil
	}

	n.log.Info().Str("request", req.RequestID).Strs("peers", blockless.PeerIDsToStr(req.Peers)).Str("consensus", req.Consensus.String()).Msg("received request to form consensus cluster")

	switch req.Consensus {
	case consensus.Raft:
		return n.createRaftCluster(ctx, from, req)

	case consensus.PBFT:
		return n.createPBFTCluster(ctx, from, req)
	}

	return fmt.Errorf("invalid consensus specified (%v %s)", req.Consensus, req.Consensus.String())
}

// processFormClusterResponse will record the cluster formation response.
func (n *Node) processFormClusterResponse(ctx context.Context, from peer.ID, res response.FormCluster) error {

	n.log.Debug().Str("request", res.RequestID).Str("from", from.String()).Msg("received cluster formation response")

	key := consensusResponseKey(res.RequestID, from)
	n.consensusResponses.Set(key, res)

	return nil
}

// processDisbandCluster will start cluster shutdown command.
func (n *Node) processDisbandCluster(ctx context.Context, from peer.ID, req request.DisbandCluster) error {

	// Should never happen.
	if !n.isWorker() {
		n.log.Warn().Str("peer", from.String()).Msg("only worker nodes participate in consensus clusters")
		return nil
	}

	n.log.Info().Str("peer", from.String()).Str("request", req.RequestID).Msg("received request to disband consensus cluster")

	err := n.leaveCluster(req.RequestID, consensusClusterDisbandTimeout)
	if err != nil {
		return fmt.Errorf("could not disband cluster (request: %s): %w", req.RequestID, err)
	}

	n.log.Info().Str("peer", from.String()).Str("request", req.RequestID).Msg("left consensus cluster")

	return nil
}

func (n *Node) formCluster(ctx context.Context, requestID string, replicas []peer.ID, consensus consensus.Type) error {

	// Create cluster formation request.
	reqCluster := request.FormCluster{
		RequestID: requestID,
		Peers:     replicas,
		Consensus: consensus,
	}

	// Request execution from peers.
	err := n.sendToMany(ctx, replicas, reqCluster, true)
	if err != nil {
		return fmt.Errorf("could not send cluster formation request to peers: %w", err)
	}

	// Wait for cluster confirmation messages.
	n.log.Debug().Str("request", requestID).Msg("waiting for cluster to be formed")

	// We're willing to wait for a limited amount of time.
	clusterCtx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
	defer exCancel()

	// Wait for confirmations for cluster forming.
	bootstrapped := make(map[string]struct{})
	var rlock sync.Mutex
	var rw sync.WaitGroup
	rw.Add(len(replicas))

	// Wait on peers asynchronously.
	for _, rp := range replicas {
		rp := rp

		go func() {
			defer rw.Done()
			key := consensusResponseKey(requestID, rp)
			res, ok := n.consensusResponses.WaitFor(clusterCtx, key)
			if !ok {
				return
			}

			n.log.Info().Str("request", requestID).Str("peer", rp.String()).Msg("accounted consensus cluster response from roll called peer")

			fc := res.(response.FormCluster)
			if fc.Code != codes.OK {
				log.Warn().Str("peer", rp.String()).Msg("peer failed to join consensus cluster")
				return
			}

			rlock.Lock()
			defer rlock.Unlock()
			bootstrapped[rp.String()] = struct{}{}
		}()
	}

	// Wait for results, whatever they may be.
	rw.Wait()

	// Err if not all peers joined the cluster successfully.
	if len(bootstrapped) != len(replicas) {
		return fmt.Errorf("some peers failed to join consensus cluster (have: %d, want: %d)", len(bootstrapped), len(replicas))
	}

	return nil
}

func (n *Node) disbandCluster(requestID string, replicas []peer.ID) error {

	msgDisband := request.DisbandCluster{
		RequestID: requestID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), consensusClusterSendTimeout)
	defer cancel()

	err := n.sendToMany(ctx, replicas, msgDisband, true)
	if err != nil {
		return fmt.Errorf("could not send cluster disband request (request: %s): %w", requestID, err)
	}

	n.log.Info().Err(err).Str("request", requestID).Strs("peers", blockless.PeerIDsToStr(replicas)).Msg("sent cluster disband request")

	return nil
}

func consensusResponseKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}
