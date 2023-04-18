package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/response"
)

// TODO: peerID of the sender is a good candidate to move on to the context

type HandlerFunc func(context.Context, peer.ID, []byte) error

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Trace().
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

	n.log.Debug().
		Str("peer", from.String()).
		Str("request_id", res.RequestID).
		Msg("processing peers roll call response")

	// Check if the response is adequate.
	if res.Code != codes.Accepted {
		n.log.Info().
			Str("peer", from.String()).
			Str("code", res.Code.String()).
			Str("request_id", res.RequestID).
			Msg("skipping inadequate roll call response - unwanted code")

		return nil
	}

	// Check if we are connected to this peer.
	// Since we receive responses to roll call via direct messages - should not happen.
	connections := n.host.Network().ConnsToPeer(from)
	if len(connections) == 0 {
		n.log.Info().
			Str("peer", from.String()).
			Str("request_id", res.RequestID).
			Msg("skipping roll call response from unconnected peer")

		return nil
	}

	n.log.Info().
		Str("peer", from.String()).
		Str("request_id", res.RequestID).
		Msg("recording peers roll call response")

	// Record the response.
	n.rollCall.add(res.RequestID, res)

	return nil
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Trace().Msg("function install response received")
	return nil
}
