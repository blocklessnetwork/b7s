package node

import (
	"context"
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
)

func (n *Node) workerProcessExecute(ctx context.Context, from peer.ID, req request.Execute) error {

	metrics.IncrCounterWithLabels(functionExecutionsMetric, 1, []metrics.Label{{Name: "function", Value: req.FunctionID}})

	requestID := req.RequestID
	if requestID == "" {
		return fmt.Errorf("request ID must be set by the head node")
	}

	// TODO: attributes
	var opts []trace.SpanStartOption
	ctx, span := n.tracer.Start(ctx, spanWorkerExecute, opts...)
	defer span.End()

	log := n.log.With().Str("request", req.RequestID).Str("function", req.FunctionID).Logger()

	// NOTE: In case of an error, we do not return early from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	code, result, err := n.workerExecute(ctx, requestID, req.Timestamp, req.Request, from)
	if err != nil {
		log.Error().Err(err).Str("peer", from.String()).Msg("execution failed")
	}

	// There's little benefit to sending a response just to say we didn't execute anything.
	if code == codes.NoContent {
		log.Info().Msg("no execution done - stopping")
		return nil
	}

	metadata, err := n.cfg.MetadataProvider.Metadata(req.Request, result.Result)
	if err != nil {
		log.Error().Err(err).Msg("could not get metadata for the execution result")
	}

	log.Info().Str("code", code.String()).Msg("execution complete")

	// Cache the execution result.
	n.executeResponses.Set(requestID, result)

	// Create the execution response from the execution result.
	res := req.Response(code).WithResults(execute.ResultMap{n.host.ID(): { Result: result, Metadata: metadata})

	// Send the response, whatever it may be (success or failure).
	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

// workerExecute is called on the worker node to use its executor component to invoke the function.
func (n *Node) workerExecute(ctx context.Context, requestID string, timestamp time.Time, req execute.Request, from peer.ID) (codes.Code, execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.fstore.IsInstalled(req.FunctionID)
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

		res, err := n.executor.ExecuteFunction(ctx, requestID, req)
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
