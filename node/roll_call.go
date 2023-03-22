package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) processRollCall(ctx context.Context, from peer.ID, payload []byte) error {

	// Only workers respond to roll calls at the moment.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Debug().Msg("skipping roll call as a non-worker node")
		return nil
	}

	// Unpack the request.
	var req request.RollCall
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack request: %w", err)
	}
	req.From = from

	// Check if we have this manifest.
	functionInstalled, err := n.isFunctionInstalled(req.FunctionID)
	if err != nil {
		// We could not lookup the manifest.
		res := response.RollCall{
			Type:       blockless.MessageRollCallResponse,
			FunctionID: req.FunctionID,
			RequestID:  req.RequestID,
			Code:       response.CodeError,
		}

		sendErr := n.send(ctx, req.From, res)
		// Log send error but choose to return the original error.
		n.log.Error().
			Err(sendErr).
			Str("to", req.From.String()).
			Msg("could not send response")

		return fmt.Errorf("could not check if function is installed: %w", err)
	}

	// We don't have this function.
	if !functionInstalled {

		res := response.RollCall{
			Type:       blockless.MessageRollCallResponse,
			FunctionID: req.FunctionID,
			RequestID:  req.RequestID,
			Code:       response.CodeNotFound,
		}

		err = n.send(ctx, req.From, res)
		if err != nil {
			return fmt.Errorf("could not send response: %w", err)
		}

		// TODO: In the original code we create a function install call here.
		// However, we do it with the CID only, but the function install code
		// requires manifestURL + CID. So at the moment this code path is not
		// present here.

		return nil
	}

	// Create response.
	res := response.RollCall{
		Type:       blockless.MessageRollCallResponse,
		FunctionID: req.FunctionID,
		RequestID:  req.RequestID,
		Code:       response.CodeAccepted,
	}

	// Send message.
	err = n.send(ctx, req.From, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

// issueRollCall will create a roll call request for executing the given function.
// On successful issuance of the roll call request, we return the ID of the issued request.
func (n *Node) issueRollCall(ctx context.Context, requestID string, functionID string) error {

	// Create a roll call request.
	rollCall := request.RollCall{
		Type:       blockless.MessageRollCall,
		FunctionID: functionID,
		RequestID:  requestID,
	}

	// Publish the mssage.
	err := n.publish(ctx, rollCall)
	if err != nil {
		return fmt.Errorf("could not publish to topic: %w", err)
	}

	return nil
}
