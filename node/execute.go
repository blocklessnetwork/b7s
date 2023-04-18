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

// executeFunc is a function that handles an execution request. In case of a worker node,
// the function is executed locally. In case of a head node, a roll call request is issued,
// and the execution request is relayed to, and retrieved from, a worker node that volunteers.
// NOTE: By using `execute.Result` here as the type, if this is executed on the head node we are
// losing the information about `who` is the peer that sent us the result - the `from` field.
type executeFunc func(context.Context, string, execute.Request) (codes.Code, map[string]execute.Result, error)

func (n *Node) processExecute(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.Execute
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	requestID := req.RequestID
	if requestID == "" {
		requestID, err = newRequestID()
		if err != nil {
			return fmt.Errorf("could not generate new request ID: %w", err)
		}
	}

	// Call the appropriate function that executes the request in the appropriate way.
	// NOTE: In case of an error, we do not return from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	var execFunc executeFunc
	if n.cfg.Role == blockless.WorkerNode {
		execFunc = n.workerExecute
	} else {
		execFunc = n.headExecute
	}

	code, results, err := execFunc(ctx, requestID, createExecuteRequest(req))
	if err != nil {
		n.log.Error().
			Err(err).
			Str("peer", from.String()).
			Str("function_id", req.FunctionID).
			Msg("execution failed")
	}

	n.log.Info().
		Str("request_id", requestID).
		Int("results", len(results)).
		Str("code", code.String()).
		Msg("execution complete")

	// Cache the execution result.
	n.executeResponses.Set(requestID, results)

	// Create the execution response from the execution result.
	res := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      code,
		RequestID: requestID,
		Results:   results,
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

// workerExecute is called on the worker node to use its executor component to invoke the function.
// The return type (map) is in order to maintain the same interface as the head node - mapping the execution result to the peer that executed it.
// In this case, the peer is us.
func (n *Node) workerExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, map[string]execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.fstore.Installed(req.FunctionID)
	if err != nil {
		return codes.Error, nil, fmt.Errorf("could not lookup function in store: %w", err)
	}

	out := make(map[string]execute.Result)

	if !functionInstalled {
		return codes.NotFound, out, nil
	}

	res, err := n.executor.ExecuteFunction(requestID, req)
	out[n.ID()] = res

	if err != nil {
		return res.Code, out, fmt.Errorf("execution failed: %w", err)
	}

	return res.Code, out, nil
}

// headExecute is called on the head node. The head node will publish a roll call and delegate an execution request to chosen nodes.
// The returned map contains execution results, mapped to the peer IDs of peers who reported them.
func (n *Node) headExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, map[string]execute.Result, error) {

	quorum := 1
	if req.Config.NodeCount > 1 {
		quorum = req.Config.NodeCount
	}

	n.log.Info().
		Str("request_id", requestID).
		Int("quorum", quorum).
		Msg("processing execution request")

	err := n.issueRollCall(ctx, requestID, req.FunctionID)
	if err != nil {
		return codes.Error, nil, fmt.Errorf("could not issue roll call: %w", err)
	}

	n.log.Info().
		Str("function_id", req.FunctionID).
		Str("request_id", requestID).
		Msg("roll call published")

	// Limit for how long we wait for responses.
	tctx, cancel := context.WithTimeout(ctx, n.cfg.RollCallTimeout)
	defer cancel()

	// Peers that have reported on roll call.
	var reportingPeers []peer.ID
rollCallResponseLoop:
	for {
		// Wait for responses from nodes who want to work on the request.
		select {
		// Request timed out.
		case <-tctx.Done():

			n.log.Warn().
				Str("function_id", req.FunctionID).
				Str("request_id", requestID).
				Msg("roll call timed out")

			return codes.Timeout, nil, blockless.ErrRollCallTimeout

		case reply := <-n.rollCall.responses(requestID):

			// Check if this is the reply we want - shouldn't really happen.
			if reply.FunctionID != req.FunctionID {

				n.log.Debug().
					Str("peer", reply.From.String()).
					Str("request_id", requestID).
					Str("function_got", reply.FunctionID).
					Str("function_want", req.FunctionID).
					Msg("skipping inadequate roll call response - wrong function")

				continue
			}

			n.log.Info().
				Str("request_id", requestID).
				Str("peer", reply.From.String()).
				Int("want_peers", quorum).
				Msg("roll called peer chosen for execution")

			reportingPeers = append(reportingPeers, reply.From)

			if len(reportingPeers) >= quorum {
				n.log.Info().Str("request_id", requestID).Int("want", quorum).Msg("enough peers reported for roll call")
				break rollCallResponseLoop
			}
		}
	}

	peerIDs := make([]string, 0, len(reportingPeers))
	for _, rp := range reportingPeers {
		peerIDs = append(peerIDs, rp.String())
	}

	n.log.Info().
		Strs("peers", peerIDs).
		Str("function_id", req.FunctionID).
		Str("request_id", requestID).
		Msg("requesting execution from peers who reported for roll call")

	// Create execution request.
	reqExecute := request.Execute{
		Type:       blockless.MessageExecute,
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
		RequestID:  requestID,
	}

	// Request execution from peers.
	for _, rp := range reportingPeers {

		err = n.send(ctx, rp, reqExecute)
		if err != nil {

			return codes.Error, nil, fmt.Errorf("could not send execution request to peer (peer: %s, function: %s, request: %s): %w",
				rp.String(),
				req.FunctionID,
				requestID,
				err)
		}
	}

	n.log.Debug().
		Int("want", quorum).
		Str("request_id", requestID).
		Msg("waiting for execution responses")

	// we're willing to wait for a limited amount of time.
	exctx, cancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
	defer cancel()

	// Wait for multiple executions.
	results := make(map[string]execute.Result)
	var rlock sync.Mutex
	var rw sync.WaitGroup
	rw.Add(len(reportingPeers))

	// Wait on peers asynchronously.
	for _, rp := range reportingPeers {
		rp := rp

		go func() {
			defer rw.Done()
			key := executionResultKey(requestID, rp)
			res, ok := n.executeResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			n.log.Debug().
				Str("request_id", requestID).
				Str("peer", rp.String()).
				Msg("accounted execution response from roll called peer")

			er := res.(response.Execute)
			// Check if there's an actual result there.
			exres, ok := er.Results[rp.String()]
			if !ok {
				return
			}

			rlock.Lock()
			defer rlock.Unlock()
			results[rp.String()] = exres
		}()
	}

	// Wait for results, whatever they may be.
	rw.Wait()

	if len(results) != quorum {
		n.log.Warn().
			Str("request_id", requestID).
			Int("have", len(results)).
			Int("want", quorum).
			Msg("did not receive enough execution responses")

		return codes.Error, nil, blockless.ErrExecutionNotEnoughNodes
	}

	n.log.Info().
		Str("request_id", requestID).
		Msg("received enough execution responses")

	code := determineOverallCode(results)

	return code, results, nil
}

func (n *Node) processExecuteResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.Execute
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not not unpack execute response: %w", err)
	}
	res.From = from

	n.log.Debug().
		Str("request_id", res.RequestID).
		Str("from", from.String()).
		Msg("received execution response")

	key := executionResultKey(res.RequestID, from)
	n.executeResponses.Set(key, res)

	return nil
}

func executionResultKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}

// determineOverallCode will return the resulting code from a set of results. Rules are:
// - if there's a single result, we use that results code
// - return OK if at least one result was successful
// - return error if none of the results were successful
func determineOverallCode(results map[string]execute.Result) codes.Code {

	if len(results) == 0 {
		return codes.NoContent
	}

	// For a single peer, just return its code.
	if len(results) == 1 {
		for peer := range results {
			return results[peer].Code
		}
	}

	// For multiple results - return OK if any of them succeeded.
	for _, res := range results {
		if res.Code == codes.OK {
			return codes.OK
		}
	}

	return codes.Error
}
