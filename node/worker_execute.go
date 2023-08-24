package node

import (
	"context"
	"encoding/json"
	"fmt"

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

	log := n.log.With().Str("request", req.RequestID).Str("function", req.FunctionID).Logger()

	// NOTE: In case of an error, we do not return early from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	code, result, err := n.workerExecute(ctx, requestID, createExecuteRequest(req), req.From)
	if err != nil {
		log.Error().Err(err).Str("peer", from.String()).Msg("execution failed")
	}

	// There's little benefit to sending a response just to say we didn't execute anything.
	if code == codes.NoContent {
		log.Info().Msg("no execution done - stopping")
		return nil
	}

	log.Info().Str("code", code.String()).Msg("execution complete")

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
	// TODO: Use the request for this, not the cluster check.
	n.clusterLock.RLock()
	cluster, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	// We are not part of a cluster - just execute the request.
	if !ok {
		res, err := n.executor.ExecuteFunction(requestID, req)
		if err != nil {
			return res.Code, res, fmt.Errorf("execution failed: %w", err)
		}

		return res.Code, res, nil
	}

	log := n.log.With().Str("request", requestID).Str("function", req.FunctionID).Logger()

	log.Info().Msg("execution request to be executed as part of a cluster")

	code, value, err := cluster.Execute(from, requestID, req)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("execution failed: %w", err)
	}

	log.Info().Str("code", string(code)).Msg("node processed the execution request")

	return code, value, nil
}
