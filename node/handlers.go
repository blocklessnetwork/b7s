package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

// TODO: peerID of the sender is a good candidate to move on to the context

type HandlerFunc func(context.Context, peer.ID, []byte) error

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Trace().Str("from", from.String()).Msg("peer health check received")
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

	log := n.log.With().Str("request", res.RequestID).Str("peer", from.String()).Logger()

	log.Debug().Msg("processing peers roll call response")

	// Check if the response is adequate.
	if res.Code != codes.Accepted {
		log.Info().Str("code", res.Code.String()).Msg("skipping inadequate roll call response - unwanted code")
		return nil
	}

	// Check if there's an active roll call already.
	exists := n.rollCall.exists(res.RequestID)
	if !exists {
		log.Info().Msg("no pending roll call for the given request, dropping response")
		return nil
	}

	log.Info().Msg("recording roll call response")

	// Record the response.
	n.rollCall.add(res.RequestID, res)

	return nil
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Trace().Str("from", from.String()).Msg("function install response received")
	return nil
}
