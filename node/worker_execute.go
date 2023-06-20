package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/raft"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// workerExecute is called on the worker node to use its executor component to invoke the function.
func (n *Node) workerExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.fstore.Installed(req.FunctionID)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not lookup function in store: %w", err)
	}

	if !functionInstalled {
		return codes.NotFound, execute.Result{}, nil
	}

	// Determine if we should just execute this function, or are we part of the cluster.
	n.clusterLock.RLock()
	raftNode, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	// There's no cluster handle created - it means we only got a direct execution request.
	if !ok {
		res, err := n.executor.ExecuteFunction(requestID, req)
		if err != nil {
			return res.Code, res, fmt.Errorf("execution failed: %w", err)
		}

		return res.Code, res, nil
	}

	// We're a part of a cluster - for now acknowledge it and return an error.
	n.log.Info().Str("request_id", requestID).Msg("execution request to be executed as part of a cluster")

	if raftNode.State() != raft.Leader {
		_, id := raftNode.LeaderWithID()

		n.log.Info().Str("request_id", requestID).Str("leader", string(id)).Msg("we're not the cluster leader - dropping the request")
		// TODO: (raft) see how we should handle this scenario upwards.
		return codes.NoContent, execute.Result{}, nil
	}

	n.log.Info().Str("request_id", requestID).Msg("we're the cluster leader, executing the request")

	fsmReq := fsmLogEntry{
		RequestID: requestID,
		Execute:   req,
	}

	payload, err := json.Marshal(fsmReq)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not serialize request for FSM")
	}

	// Apply Raft log.
	future := raftNode.Apply(payload, defaultRaftApplyTimeout)
	err = future.Error()
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not apply raft log: %w", err)
	}

	n.log.Info().Str("request_id", requestID).Msg("node applied raft log")

	// Get execution result.
	response := future.Response()
	value, ok := response.(execute.Result)
	if !ok {
		fsmErr, ok := response.(error)
		if ok {
			return codes.Error, execute.Result{}, fmt.Errorf("execution encountered an error: %w", fsmErr)
		}

		n.log.Error().Type("type", response).Msg("unexpected FSM response format")

		return codes.Error, execute.Result{}, errors.New("unexpected FSM response format")
	}

	n.log.Info().Str("request_id", requestID).Msg("cluster leader executed the request")

	return codes.OK, value, nil
}
