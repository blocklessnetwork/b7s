package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/request"
)

func (n *Node) processExecute(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request.
	var req request.Execute
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack the request: %w", err)
	}

	req.From = from

	// TODO: Sent to local channel executor - should be handled here.
	return fmt.Errorf("TBD: Not completed")
}

func (n *Node) processExecuteResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the request
	var req request.ExecuteResponse
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not not unpack the request: %w", err)
	}

	// TODO: Complete this flow.

	return errors.New("TBD: Not implemented")
}
