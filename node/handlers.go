package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO: Check - what do these functions use contexts for?
// TODO: peerID of the sender is a good candidate to move on to the context

type HandlerFunc func(context.Context, peer.ID, []byte) error

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Debug().Msg("peer health check received")
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
	n.recordRollCallResponse(res)

	return nil
}

func (n *Node) recordRollCallResponse(res response.RollCall) {
	n.rollCallResponses[res.RequestID] <- res
}

func (n *Node) processInstallFunction(ctx context.Context, from peer.ID, payload []byte) error {
	return errors.New("TBD: Not implemented")
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Debug().Msg("function install response received")
	return nil
}
