package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) processFormCluster(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.FormCluster
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	n.log.Info().Str("request_id", req.RequestID).Strs("peers", peerIDList(req.Peers)).Msg("received request to form consensus cluster")

	raftHandler, err := n.newRaftHandler(req.RequestID)
	if err != nil {
		return fmt.Errorf("could not create raft node: %w", err)
	}

	err = bootstrapCluster(raftHandler, req.Peers)
	if err != nil {
		return fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	n.clusterLock.Lock()
	n.clusters[req.RequestID] = raftHandler
	n.clusterLock.Unlock()

	n.log.Info().Msg("waiting on leadership notification")

	// Wait until we have leadership info to confirm.
	isLeader := <-raftHandler.LeaderCh()

	n.log.Info().Bool("leader", isLeader).Msg("notified of leadership change")

	res := response.FormCluster{
		Type:      blockless.MessageFormClusterResponse,
		RequestID: req.RequestID,
		Code:      codes.OK,
	}

	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send cluster confirmation message: %w", err)
	}

	return nil
}

// processFormClusterResponse will record the cluster formation response
func (n *Node) processFormClusterResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.FormCluster
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not unpack the cluster formation response: %w", err)
	}
	res.From = from

	n.log.Debug().Str("request_id", res.RequestID).Str("from", from.String()).Msg("received cluster formation response")

	key := consensusResponseKey(res.RequestID, from)
	n.consensusResponses.Set(key, res)

	return nil
}

func consensusResponseKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}
