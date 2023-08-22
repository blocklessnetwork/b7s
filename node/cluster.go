package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/consensus"
	"github.com/blocklessnetworking/b7s/consensus/raft"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) processFormCluster(ctx context.Context, from peer.ID, payload []byte) error {

	// Should never happen.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Warn().Str("peer", from.String()).Msg("only worker nodes participate in consensus clusters")
		return nil
	}

	// Unpack the request.
	var req request.FormCluster
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	n.log.Info().Str("request", req.RequestID).Strs("peers", blockless.PeerIDsToStr(req.Peers)).Str("consensus", req.Consensus.String()).Msg("received request to form consensus cluster")

	switch req.Consensus {
	case consensus.Raft:
		return n.createRaftCluster(req)

	case consensus.PBFT:
		return fmt.Errorf("TBD: PBFT implementation coming soon")
	}

	return fmt.Errorf("invalid consensus specified (%s %s)", req.Consensus, req.Consensus.String())
}

// processFormClusterResponse will record the cluster formation response.
func (n *Node) processFormClusterResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.FormCluster
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not unpack the cluster formation response: %w", err)
	}
	res.From = from

	n.log.Debug().Str("request", res.RequestID).Str("from", from.String()).Msg("received cluster formation response")

	key := consensusResponseKey(res.RequestID, from)
	n.consensusResponses.Set(key, res)

	return nil
}

// processDisbandCluster will start cluster shutdown command.
func (n *Node) processDisbandCluster(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.DisbandCluster
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	n.log.Info().Str("request", req.RequestID).Msg("received request to disband consensus cluster")

	err = n.leaveCluster(req.RequestID)
	if err != nil {
		return fmt.Errorf("could not disband cluster (request: %s): %w", req.RequestID, err)
	}

	return nil
}

func (n *Node) leaveCluster(requestID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), raftClusterDisbandTimeout)
	defer cancel()

	// We know that the request is done executing when we have a result for it.
	_, ok := n.executeResponses.WaitFor(ctx, requestID)

	n.log.Info().Bool("executed_work", ok).Str("request", requestID).Msg("waiting for execution done, leaving raft cluster")

	log := n.log.With().Str("request", requestID).Logger()
	log.Info().Msg("shutting down cluster")

	n.clusterLock.RLock()
	raftHandler, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	if !ok {
		return nil
	}

	err := raftHandler.Shutdown()
	if err != nil {
		return fmt.Errorf("could not leave raft cluster (request: %v): %w", requestID, err)
	}

	n.clusterLock.Lock()
	delete(n.clusters, requestID)
	n.clusterLock.Unlock()

	return nil
}

func consensusResponseKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}

func (n *Node) createRaftCluster(fc request.FormCluster) error {

	// Add a callback function to cache the execution result
	cacheFn := func(req raft.FSMLogEntry, res execute.Result) {
		n.executeResponses.Set(req.RequestID, res)
	}

	// Add a callback function to send the execution result to origin.
	sendFn := func(req raft.FSMLogEntry, res execute.Result) {

		ctx, cancel := context.WithTimeout(context.Background(), raftClusterSendTimeout)
		defer cancel()

		msg := response.Execute{
			Type:      blockless.MessageExecuteResponse,
			Code:      res.Code,
			RequestID: req.RequestID,
			Results: execute.ResultMap{
				n.host.ID(): res,
			},
		}

		err := n.send(ctx, req.Origin, msg)
		if err != nil {
			n.log.Error().Err(err).Str("peer", req.Origin.String()).Msg("could not send execution result to node")
		}
	}

	rh, err := raft.New(
		n.log,
		n.host,
		n.cfg.Workspace,
		fc.RequestID,
		n.executor,
		fc.Peers,
		raft.WithCallbacks(cacheFn, sendFn),
	)
	if err != nil {
		return fmt.Errorf("could not create raft node: %w", err)
	}

	n.clusterLock.Lock()
	n.clusters[fc.RequestID] = rh
	n.clusterLock.Unlock()

	return nil
}
