package node

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, _ response.Health) error {
	n.log.Trace().Str("from", from.String()).Msg("peer health check received")
	return nil
}

func (n *Node) processRollCallResponse(ctx context.Context, from peer.ID, res response.RollCall) error {

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

	rres := rollCallResponse{
		From:     from,
		RollCall: res,
	}

	// Record the response.
	n.rollCall.add(res.RequestID, rres)

	return nil
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, _ response.InstallFunction) error {
	n.log.Trace().Str("from", from.String()).Msg("function install response received")
	return nil
}
