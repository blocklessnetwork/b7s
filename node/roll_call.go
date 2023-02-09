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
	if n.role != blockless.WorkerNode {
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
	_, err = n.getFunctionManifest(req.FunctionID)
	if err != nil {

		// TODO: Install this function now.

		// Notify the caller that we don't have this manifest.

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
