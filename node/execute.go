package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
type executeFunc func(context.Context, string, execute.Request) (codes.Code, execute.Result, error)

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

	execFunc := n.getExecuteFunction(req)

	// Call the appropriate function that executes the request in the appropriate way.
	// NOTE: In case of an error, we do not return early from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	code, result, err := execFunc(ctx, requestID, createExecuteRequest(req))
	if err != nil {
		n.log.Error().
			Err(err).
			Str("peer", from.String()).
			Str("function_id", req.FunctionID).
			Msg("execution failed")
	}

	// There's little benefit to sending a response just to say we didn't execute anything.
	if code == codes.NoContent {
		n.log.Info().Str("request_id", requestID).Msg("no execution done - stopping")
		return nil
	}

	n.log.Info().
		Str("request_id", requestID).
		Str("code", code.String()).
		Msg("execution complete")

	// Cache the execution result.
	n.executeResponses.Set(requestID, result)

	// Create the execution response from the execution result.
	res := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      code,
		RequestID: requestID,
		Result:    result,
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

// helper function to to convert a slice of multiaddrs to strings
func peerIDList(ids []peer.ID) []string {
	peerIDs := make([]string, 0, len(ids))
	for _, rp := range ids {
		peerIDs = append(peerIDs, rp.String())
	}
	return peerIDs
}

func (n *Node) getExecuteFunction(req request.Execute) executeFunc {

	if n.cfg.Role == blockless.HeadNode {
		return n.headExecute
	}

	return n.workerExecute
}
