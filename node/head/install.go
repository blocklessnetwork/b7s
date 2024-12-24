package head

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/response"
)

func (h *HeadNode) processInstallFunctionResponse(ctx context.Context, from peer.ID, res response.InstallFunction) error {

	h.Log().Trace().
		Stringer("from", from).
		Str("cid", res.CID).
		Msg("function install response received")

	return nil
}
