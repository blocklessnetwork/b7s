package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

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

	// If we're a worker node - execute the function locally.
	if n.role == blockless.WorkerNode {

		// TODO: Check if function is installed.

		// Execute the function.
		res, err := n.execute.Function(execReq)
		if err != nil {
			n.log.Error().Err(err).Msg("execution failed")
		}

		// Cache the execution result.
		n.excache.Set(res.RequestID, &res)

		// Send the response, whatever it may be (success or failure).
		err = n.send(ctx, req.From, res)
		if err != nil {
			return fmt.Errorf("could not send response: %w", err)
		}
	}

	return n.headNodeExecute(ctx, from, execReq)
}

func (n *Node) headNodeExecute(ctx context.Context, from peer.ID, req execute.Request) error {

	requestID, err := n.issueRollCall(ctx, req.FunctionID)
	if err != nil {
		return fmt.Errorf("could not issue roll call: %w", err)
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

			res := execute.Response{
				Code: response.CodeTimeout,
			}
			_ = res

			// TODO: Who do we send to?

		case reply := <-n.rollCallResponses[requestID]:

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

	// Create a channel where execution response will be received.
	n.executeResponses[requestID] = make(chan response.Execute)

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
		// TODO: Send response to caller.
		return fmt.Errorf("could not send execution request to peer (peer: %s, function: %s, request: %s): %w",
			reportingPeer.String(),
			req.FunctionID,
			requestID,
			err)
	}

	resExecute := <-n.executeResponses[requestID]

	n.log.Info().
		Str("request_id", requestID).
		Str("peer", resExecute.From.String()).
		Str("code", resExecute.Code).
		Msg("received execution response")

	// Return the execution response.
	// TODO: Use interfaces for worker and head node - for execution handlers.

	out := execute.Response{
		Code:      resExecute.Code,
		Result:    resExecute.Result,
		RequestID: resExecute.RequestID,
	}

	// TODO: Execution cache.

	err = n.send(ctx, from, out)
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

	// Record execution response.
	n.recordExecuteResponse(res)

	return nil
}

func (n *Node) recordExecuteResponse(res response.Execute) {
	n.executeResponses[res.RequestID] <- res
}
