package head

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/response"
)

func (h *HeadNode) processHealthCheck(ctx context.Context, from peer.ID, _ response.Health) error {
	h.Log().Trace().Stringer("from", from).Msg("peer health check received")
	return nil
}
