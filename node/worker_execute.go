package node

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
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
	code, result, err := n.workerExecute(ctx, requestID, req.Timestamp, createExecuteRequest(req), req.From)
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
func (n *Node) workerExecute(ctx context.Context, requestID string, timestamp time.Time, req execute.Request, from peer.ID) (codes.Code, execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.fstore.Installed(req.FunctionID)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not lookup function in store: %w", err)
	}

	if !functionInstalled {
		return codes.NotFound, execute.Result{}, nil
	}

	// Determine if we should just execute this function, or are we part of the cluster.

	// Here we actually have a bit of a conceptual problem with having the same models for head and worker node.
	// Head node receives client requests so it can expect _some_ type of inaccuracy there. Worker node receives
	// execution requests from the head node, so it shouldn't really tolerate errors/ambiguities.
	consensus, err := parseConsensusAlgorithm(req.Config.ConsensusAlgorithm)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not parse consensus algorithm from the head node request, aborting (value: %s): %w", req.Config.ConsensusAlgorithm, err)
	}

	// We are not part of a cluster - just execute the request.
	if !consensusRequired(consensus) {
		res, err := n.executor.ExecuteFunction(requestID, req)
		if err != nil {
			return res.Code, res, fmt.Errorf("execution failed: %w", err)
		}

		return res.Code, res, nil
	}

	// Now we KNOW we need a consensus. A cluster must already exist.

	n.clusterLock.RLock()
	cluster, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	if !ok {
		return codes.Error, execute.Result{}, fmt.Errorf("consensus required but no cluster found; omitted cluster formation message or error forming cluster (request: %s)", requestID)
	}

	log := n.log.With().Str("request", requestID).Str("function", req.FunctionID).Str("consensus", consensus.String()).Logger()

	log.Info().Msg("execution request to be executed as part of a cluster")

	code, value, err := cluster.Execute(from, requestID, timestamp, req)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("execution failed: %w", err)
	}

	log.Info().Str("code", string(code)).Msg("node processed the execution request")

	return code, value, nil
}
