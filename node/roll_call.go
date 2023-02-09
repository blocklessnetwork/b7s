package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/request"
)

func (n *Node) processRollCall(ctx context.Context, from peer.ID, payload []byte) error {

	// Only workers should respond to roll-calls.
	// NOTE: Potentially other nodes should be able to respond to roll calls too.
	if n.role != blockless.WorkerNode {
		n.log.Debug().Msg("skipping roll-call as a non-worker node")
		return nil
	}

	// Unpack the request.
	var req request.RollCall
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack request: %w", err)
	}
	req.From = from

	// TODO: Complete this flow.

	return errors.New("TBD: Not implemented")
}
