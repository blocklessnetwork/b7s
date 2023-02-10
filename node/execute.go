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

	// TODO: Head node implement.
	return nil
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
