package head

import (
	"context"

	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (h *HeadNode) processInstallFunctionResponse(ctx context.Context, from peer.ID, res response.InstallFunction) error {

	h.Log().Trace().
		Stringer("from", from).
		Str("cid", res.CID).
		Msg("function install response received")

	return nil
}
