package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) workerProcessExecute(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.Execute
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	requestID := req.RequestID
	if requestID == "" {
		return fmt.Errorf("request ID must be set by the head node")
	}

	// NOTE: In case of an error, we do not return early from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	code, result, err := n.workerExecute(ctx, requestID, createExecuteRequest(req), req.From)
	if err != nil {
		n.log.Error().Err(err).Str("peer", from.String()).Str("function_id", req.FunctionID).Str("request_id", requestID).Msg("execution failed")
	}

	// There's little benefit to sending a response just to say we didn't execute anything.
	if code == codes.NoContent {
		n.log.Info().Str("request_id", requestID).Msg("no execution done - stopping")
		return nil
	}

	n.log.Info().Str("request_id", requestID).Str("code", code.String()).Msg("execution complete")

	// Cache the execution result.
	n.executeResponses.Set(requestID, result)

	// Create the execution response from the execution result.
	res := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      code,
		RequestID: requestID,
		Results: execute.ResultMap{
			n.host.ID(): result,
		},
	}

	// Send the response, whatever it may be (success or failure).
	err = n.send(ctx, req.From, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

// workerExecute is called on the worker node to use its executor component to invoke the function.
func (n *Node) workerExecute(ctx context.Context, requestID string, req execute.Request, from peer.ID) (codes.Code, execute.Result, error) {

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

	// We are not part of a cluster - just execute the request.
	if !ok {
		res, err := n.executor.ExecuteFunction(requestID, req)
		if err != nil {
			return res.Code, res, fmt.Errorf("execution failed: %w", err)
		}

		return res.Code, res, nil
	}

	n.log.Info().Str("request_id", requestID).Msg("execution request to be executed as part of a cluster")

	if raftNode.State() != raft.Leader {
		_, id := raftNode.LeaderWithID()

		n.log.Info().Str("request_id", requestID).Str("leader", string(id)).Msg("we are not the cluster leader - dropping the request")
		return codes.NoContent, execute.Result{}, nil
	}

	n.log.Info().Str("request_id", requestID).Msg("we are the cluster leader, executing the request")

	fsmReq := fsmLogEntry{
		RequestID: requestID,
		Origin:    from,
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

		return codes.Error, execute.Result{}, fmt.Errorf("unexpected FSM response format: %T", response)
	}

	n.log.Info().Str("request_id", requestID).Msg("cluster leader executed the request")

	return codes.OK, value, nil
}
