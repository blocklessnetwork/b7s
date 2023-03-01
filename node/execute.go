package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

// executeFunc is a function that handles an execution request. In case of a worker node,
// the function is executed locally. In case of a head node, a roll call request is issued,
// and the execution request is relayed to, and retrieved from, a worker node that volunteers.
// NOTE: By using `execute.Result` here as the type, if this is executed on the head node we are
// losing the information about `who` is the peer that sent us the result - the `from` field.
type executeFunc func(context.Context, peer.ID, execute.Request) (execute.Result, error)

func (n *Node) processExecute(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.Execute
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}
	req.From = from

	// Create execute request.
	execReq := execute.Request{
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
	}

	// Call the appropriate function that executes the request in the appropriate way.
	// NOTE: In case of an error, we do not return from this function.
	// Instead, we send the response back to the caller, whatever it may be.
	var execFunc executeFunc
	if n.role == blockless.WorkerNode {
		execFunc = n.workerExecute
	} else {
		execFunc = n.headExecute
	}

	result, err := execFunc(ctx, from, execReq)
	if err != nil {
		n.log.Error().
			Err(err).
			Str("peer", from.String()).
			Str("function_id", req.FunctionID).
			Msg("execution failed")
	}

	// Cache the execution result.
	n.excache.Set(result.RequestID, result)

	// Create the execution response from the execution result.
	res := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		RequestID: result.RequestID,
		Code:      result.Code,
		Result:    result.Result,
	}

	// Send the response, whatever it may be (success or failure).
	err = n.send(ctx, req.From, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

func (n *Node) workerExecute(ctx context.Context, from peer.ID, req execute.Request) (execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.isFunctionInstalled(req.FunctionID)
	if err != nil {
		res := execute.Result{
			Code: response.CodeError,
		}
		return res, fmt.Errorf("could not lookup function in store: %w", err)
	}

	if !functionInstalled {
		res := execute.Result{
			Code: response.CodeNotFound,
		}

		return res, nil
	}

	// Execute the function.
	res, err := n.execute.Function(req)
	if err != nil {
		return res, fmt.Errorf("execution failed: %w", err)
	}

	return res, nil
}

func (n *Node) headExecute(ctx context.Context, from peer.ID, req execute.Request) (execute.Result, error) {

	requestID, err := n.issueRollCall(ctx, req.FunctionID)
	if err != nil {

		res := execute.Result{
			Code: response.CodeError,
		}

		return res, fmt.Errorf("could not issue roll call: %w", err)
	}

	n.log.Info().
		Str("function_id", req.FunctionID).
		Str("request_id", requestID).
		Msg("roll call published")

	// Limit for how long we wait for responses.
	tctx, cancel := context.WithTimeout(ctx, rollCallTimeout)
	defer cancel()

	// Peer that reports to roll call first.
	var reportingPeer peer.ID
rollCallResponseLoop:
	for {
		// Wait for responses from nodes who want to work on the request.
		select {
		// Request timed out.
		case <-tctx.Done():

			n.log.Info().
				Str("function_id", req.FunctionID).
				Str("request_id", requestID).
				Msg("roll call timed out")

			res := execute.Result{
				Code: response.CodeTimeout,
			}

			return res, errRollCallTimeout

		case reply := <-n.rollCall.responses(requestID):

			n.log.Debug().
				Str("peer", reply.From.String()).
				Str("function_id", req.FunctionID).
				Str("request_id", requestID).
				Msg("peer reported for roll call")

			// Check if this is the reply we want.
			if reply.Code != response.CodeAccepted ||
				reply.FunctionID != req.FunctionID ||
				reply.RequestID != requestID {
				continue
			}

			// Check if we are connected to this peer.
			connections := n.host.Network().ConnsToPeer(reply.From)
			if len(connections) == 0 {
				continue
			}

			reportingPeer = reply.From
			break rollCallResponseLoop
		}
	}

	n.log.Info().
		Str("peer", reportingPeer.String()).
		Str("function_id", req.FunctionID).
		Str("request_id", requestID).
		Msg("peer reported for roll call")

	// Request execution from the peer who reported back first.
	reqExecute := request.Execute{
		Type:       blockless.MessageExecute,
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
	}

	// Send message to reporting peer to execute the function.
	err = n.send(ctx, reportingPeer, reqExecute)
	if err != nil {

		res := execute.Result{
			Code: response.CodeError,
		}

		return res, fmt.Errorf("could not send execution request to peer (peer: %s, function: %s, request: %s): %w",
			reportingPeer.String(),
			req.FunctionID,
			requestID,
			err)
	}

	// TODO: Verify that the response came from the peer that reported for the roll call.
	resExecute := n.executeResponses.Wait(requestID).(response.Execute)

	n.log.Info().
		Str("request_id", requestID).
		Str("peer", resExecute.From.String()).
		Str("code", resExecute.Code).
		Msg("received execution response")

	// Return the execution result.
	result := execute.Result{
		Code:      resExecute.Code,
		Result:    resExecute.Result,
		RequestID: resExecute.RequestID,
	}

	return result, nil
}

func (n *Node) processExecuteResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.Execute
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not not unpack execute response: %w", err)
	}
	res.From = from

	// Record execution response.
	n.executeResponses.Set(res.RequestID, res)

	return nil
}

// isFuncitonInstalled looks up the function in the store by using the functionID/CID as key.
func (n *Node) isFunctionInstalled(functionID string) (bool, error) {

	_, err := n.function.Get("", functionID, true)
	if err != nil {

		if errors.Is(err, blockless.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("could not lookup function in store: %w", err)
	}

	return true, nil
}
