package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/raft"
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

	// TODO: (raft) Stopgap, have this done correctly.
	// Register an observer to monitor leadership changes. More precisely,
	// wait on the first leader election, so we know when the cluster is operational.

	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, func(obs *raft.Observation) bool {
		_, ok := obs.Data.(raft.LeaderObservation)
		return ok
	})

	// TODO: (raft) - consider making this synchronous
	go func() {
		// Wait on leadership observation.
		obs := <-obsCh
		leaderObs, ok := obs.Data.(raft.LeaderObservation)
		if !ok {
			n.log.Error().Type("type", obs.Data).Msg("invalid observation type received")
			return
		}

		// We don't need the observer anymore.
		raftHandler.DeregisterObserver(observer)

		n.log.Info().Str("peer", from.String()).Str("leader", string(leaderObs.LeaderID)).Msg("observed a leadership event - sending response")

		res := response.FormCluster{
			Type:      blockless.MessageFormClusterResponse,
			RequestID: req.RequestID,
			Code:      codes.OK,
		}

		err = n.send(ctx, from, res)
		if err != nil {
			n.log.Error().Err(err).Msg("could not send cluster confirmation message")
			return
		}
	}()

	raftHandler.RegisterObserver(observer)

	err = bootstrapCluster(raftHandler, req.Peers)
	if err != nil {
		return fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	n.clusterLock.Lock()
	n.clusters[req.RequestID] = raftHandler
	n.clusterLock.Unlock()

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