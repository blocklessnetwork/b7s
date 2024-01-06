package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

// TODO: Check - head node really accepts execution requests from the REST API. Should this message handling be cognizant of `topics`?
func (n *Node) headProcessExecute(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.Execute
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	requestID, err := newRequestID()
	if err != nil {
		return fmt.Errorf("could not generate new request ID: %w", err)
	}

	log := n.log.With().Str("request", req.RequestID).Str("peer", from.String()).Str("function", req.FunctionID).Logger()

	code, results, cluster, err := n.headExecute(ctx, requestID, req.Request, "")
	if err != nil {
		log.Error().Err(err).Msg("execution failed")
	}

	log.Info().Str("code", code.String()).Msg("execution complete")

	// Create the execution response from the execution result.
	res := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      code,
		RequestID: requestID,
		Results:   results,
		Cluster:   cluster,
	}

	// Communicate the reason for failure in these cases.
	if errors.Is(err, blockless.ErrRollCallTimeout) || errors.Is(err, blockless.ErrExecutionNotEnoughNodes) {
		res.Message = err.Error()
	}

	// Send the response, whatever it may be (success or failure).
	err = n.send(ctx, req.From, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

// headExecute is called on the head node. The head node will publish a roll call and delegate an execution request to chosen nodes.
// The returned map contains execution results, mapped to the peer IDs of peers who reported them.
func (n *Node) headExecute(ctx context.Context, requestID string, req execute.Request, subgroup string) (codes.Code, execute.ResultMap, execute.Cluster, error) {

	nodeCount := 1
	if req.Config.NodeCount > 1 {
		nodeCount = req.Config.NodeCount
	}

	// Create a logger with relevant context.
	log := n.log.With().Str("request", requestID).Str("function", req.FunctionID).Int("node_count", nodeCount).Logger()

	consensusAlgo, err := parseConsensusAlgorithm(req.Config.ConsensusAlgorithm)
	if err != nil {
		log.Error().Str("value", req.Config.ConsensusAlgorithm).Str("default", n.cfg.DefaultConsensus.String()).Err(err).Msg("could not parse consensus algorithm from the user request, using default")
		consensusAlgo = n.cfg.DefaultConsensus
	}

	if consensusRequired(consensusAlgo) {
		log = log.With().Str("consensus", consensusAlgo.String()).Logger()
	}

	log.Info().Msg("processing execution request")

	// Phase 1. - Issue roll call to nodes.
	reportingPeers, err := n.executeRollCall(ctx, requestID, req.FunctionID, nodeCount, consensusAlgo, subgroup, req.Config.Attributes)
	if err != nil {
		code := codes.Error
		if errors.Is(err, blockless.ErrRollCallTimeout) {
			code = codes.Timeout
		}

		return code, nil, execute.Cluster{}, fmt.Errorf("could not roll call peers (request: %s): %w", requestID, err)
	}

	cluster := execute.Cluster{
		Peers: reportingPeers,
	}

	// Phase 2. - Request cluster formation, if we need consensus.
	if consensusRequired(consensusAlgo) {

		log.Info().Strs("peers", blockless.PeerIDsToStr(reportingPeers)).Msg("requesting cluster formation from peers who reported for roll call")

		err := n.formCluster(ctx, requestID, reportingPeers, consensusAlgo)
		if err != nil {
			return codes.Error, nil, execute.Cluster{}, fmt.Errorf("could not form cluster (request: %s): %w", requestID, err)
		}

		// When we're done, send a message to disband the cluster.
		// NOTE: We could schedule this on the worker nodes when receiving the execution request.
		// One variant I tried is waiting on the execution to be done on the leader (using a timed wait on the execution response) and starting raft shutdown after.
		// However, this can happen too fast and the execution request might not have been propagated to all of the nodes in the cluster, but "only" to a majority.
		// Doing this here allows for more wiggle room and ~probably~ all nodes will have seen the request so far.
		defer n.disbandCluster(requestID, reportingPeers)
	}

	// Phase 3. - Request execution.

	// Send the execution request to peers in the cluster. Non-leaders will drop the request.
	reqExecute := request.Execute{
		Type:      blockless.MessageExecute,
		Request:   req,
		RequestID: requestID,
		Timestamp: time.Now().UTC(),
	}

	// If we're working with PBFT, sign the request.
	if consensusAlgo == consensus.PBFT {
		err := reqExecute.Request.Sign(n.host.PrivateKey())
		if err != nil {
			return codes.Error, nil, cluster, fmt.Errorf("could not sign execution request (function: %s, request: %s): %w", req.FunctionID, requestID, err)
		}
	}

	err = n.sendToMany(ctx, reportingPeers, reqExecute)
	if err != nil {
		return codes.Error, nil, cluster, fmt.Errorf("could not send execution request to peers (function: %s, request: %s): %w", req.FunctionID, requestID, err)
	}

	log.Debug().Msg("waiting for execution responses")

	var results execute.ResultMap
	if consensusAlgo == consensus.PBFT {
		results = n.gatherExecutionResultsPBFT(ctx, requestID, reportingPeers)

		log.Info().Msg("received PBFT execution responses")

		retcode := codes.OK
		// Use the return code from the execution as the return code.
		for _, res := range results {
			retcode = res.Code
			break
		}

		return retcode, results, cluster, nil
	}

	results = n.gatherExecutionResults(ctx, requestID, reportingPeers)

	log.Info().Int("cluster_size", len(reportingPeers)).Int("responded", len(results)).Msg("received execution responses")

	// How many results do we have, and how many do we expect.
	respondRatio := float64(len(results)) / float64(len(reportingPeers))
	threshold := determineThreshold(req)

	retcode := codes.OK
	if respondRatio == 0 {
		retcode = codes.NoContent
	} else if respondRatio < threshold {
		log.Warn().Float64("expected", threshold).Float64("have", respondRatio).Msg("threshold condition not met")
		retcode = codes.PartialContent
	}

	return retcode, results, cluster, nil
}

func determineThreshold(req execute.Request) float64 {

	if req.Config.Threshold > 0 && req.Config.Threshold <= 1 {
		return req.Config.Threshold
	}

	return defaultExecutionThreshold
}
