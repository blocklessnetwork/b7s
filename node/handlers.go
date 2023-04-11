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

// TODO: peerID of the sender is a good candidate to move on to the context

type HandlerFunc func(context.Context, peer.ID, []byte) error

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Debug().
		Str("from", from.String()).
		Msg("peer health check received")
	return nil
}

func (n *Node) processRollCallResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the roll call response.
	var res response.RollCall
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not unpack the roll call response: %w", err)
	}
	res.From = from

	// Record the response.
	n.rollCall.add(res.RequestID, res)

	return nil
}

func (n *Node) processInstallFunction(ctx context.Context, from peer.ID, payload []byte) error {

	// Only workers should respond to function install requests.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Debug().
			Msg("received function install request, ignoring")
		return nil
	}

	// Unpack the request.
	var req request.InstallFunction
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack request: %w", err)
	}
	req.From = from

	// Get the function manifest.
	_, err = n.fstore.Get(req.ManifestURL, req.CID, true)
	if err != nil {
		return fmt.Errorf("could not retrieve function (manifest_url: %s, cid: %s): %w", req.ManifestURL, req.CID, err)
	}

	// Create the response.
	res := response.InstallFunction{
		Type:    blockless.MessageInstallFunctionResponse,
		Code:    response.CodeAccepted,
		Message: "installed",
	}

	// Reply to the caller.
	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send the response (peer: %s): %w", from, err)
	}

	return nil
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Debug().Msg("function install response received")
	return nil
}
