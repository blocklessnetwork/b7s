package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/consensus"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) processFormCluster(ctx context.Context, from peer.ID, payload []byte) error {

	// Should never happen.
	if !n.isWorker() {
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
		return n.createRaftCluster(ctx, from, req)

	case consensus.PBFT:
		return n.createPBFTCluster(ctx, from, req)
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

	// Should never happen.
	if !n.isWorker() {
		n.log.Warn().Str("peer", from.String()).Msg("only worker nodes participate in consensus clusters")
		return nil
	}

	// Unpack the request.
	var req request.DisbandCluster
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	n.log.Info().Str("peer", from.String()).Str("request", req.RequestID).Msg("received request to disband consensus cluster")

	err = n.leaveCluster(req.RequestID)
	if err != nil {
		return fmt.Errorf("could not disband cluster (request: %s): %w", req.RequestID, err)
	}

	return nil
}

func consensusResponseKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}
