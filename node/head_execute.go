package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

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

	code, results, cluster, err := n.headExecute(ctx, requestID, createExecuteRequest(req))
	if err != nil {
		log.Error().Err(err).Msg("execution failed")
	}

	log.Info().Str("code", code.String()).Msg("execution complete")

	// NOTE: Head node no longer caches execution results because it doesn't have one of its own.

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
func (n *Node) headExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, execute.ResultMap, execute.Cluster, error) {

	// TODO: (raft) if no cluster/consensus is required - request direct execution.
	quorum := 1
	if req.Config.NodeCount > 1 {
		quorum = req.Config.NodeCount
	}

	consensus, err := parseConsensusAlgorithm(req.Config.ConsensusAlgorithm)
	if err != nil {
		n.log.Error().Str("value", req.Config.ConsensusAlgorithm).Str("default", n.cfg.DefaultConsensus.String()).Err(err).Msg("could not parse consensus algorithm from the user request, using default")
		consensus = n.cfg.DefaultConsensus
	}

	// Create a logger with relevant context.
	log := n.log.With().Str("request", requestID).Str("function", req.FunctionID).Int("quorum", quorum).Str("consenus", consensus.String()).Logger()

	log.Info().Msg("processing execution request")

	// Phase 1. - Issue roll call to nodes.

	reportingPeers, err := n.executeRollCall(ctx, requestID, req.FunctionID, quorum, consensus)
	if err != nil {
		code := codes.Error
		if errors.Is(err, blockless.ErrRollCallTimeout) {
			code = codes.Timeout
		}

		return code, nil, execute.Cluster{}, fmt.Errorf("could not roll call peers (request: %s): %w", requestID, err)
	}

	log.Info().Strs("peers", blockless.PeerIDsToStr(reportingPeers)).Msg("requesting cluster formation from peers who reported for roll call")

	cluster := execute.Cluster{
		Peers: reportingPeers,
	}

	// Phase 2. - Request cluster formation, if we need consensus.
	if consensusRequired(consensus) {

		err := n.formCluster(ctx, requestID, reportingPeers, consensus)
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
		Type:       blockless.MessageExecute,
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
		RequestID:  requestID,
	}
	err = n.sendToMany(ctx, reportingPeers, reqExecute)
	if err != nil {
		return codes.Error, nil, cluster, fmt.Errorf("could not send execution request to peers (function: %s, request: %s): %w", req.FunctionID, requestID, err)
	}

	log.Debug().Msg("waiting for execution responses")

	// We're willing to wait for a limited amount of time.
	exctx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
	defer exCancel()

	var (
		// We're waiting for a single execution result now, as only the cluster leader will return a result.
		results execute.ResultMap = make(map[peer.ID]execute.Result)
		reslock sync.Mutex
		wg      sync.WaitGroup
	)

	wg.Add(len(reportingPeers))

	// Wait on peers asynchronously.
	for _, rp := range reportingPeers {
		rp := rp

		go func() {
			defer wg.Done()
			key := executionResultKey(requestID, rp)
			res, ok := n.executeResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			log.Info().Str("peer", rp.String()).Msg("accounted execution response from peer")

			er := res.(response.Execute)

			exres, ok := er.Results[rp]
			if !ok {
				return
			}

			reslock.Lock()
			defer reslock.Unlock()
			results[rp] = exres
		}()
	}

	wg.Wait()

	log.Info().Int("cluster_size", len(reportingPeers)).Int("responded", len(results)).Msg("received execution responses")

	// TODO: Depending on the consensus, we want to treat results differently. E.g. for PBFT we may only want f+1 response and we're good.

	// How many results do we have, and how many do we expect.
	respondRatio := float64(len(results)) / float64(len(reportingPeers))
	threshold := determineThreshold(req)

	retcode := codes.OK
	if respondRatio < threshold {
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
