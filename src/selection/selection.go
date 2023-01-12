package selection

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Network interface {
	ConnsToPeer(id peer.ID) map[peer.ID]bool
}

type Host interface {
	Network() Network
}

func SelectWorkerFromRollCall(
	ctx context.Context,
	rollcallResponse models.MsgRollCallResponse,
	request models.RequestExecute,
	rollcallRequest models.MsgRollCall,
) bool {

	h := ctx.Value("host").(Host)
	conns := h.Network().ConnsToPeer(rollcallResponse.From)

	// pop off all the responses that don't match our first found connection
	// worker with function
	// worker responsed accepted has resources and is ready to execute
	// worker knows RequestID
	if rollcallResponse.Code == enums.ResponseCodeAccepted && rollcallResponse.FunctionId == request.FunctionId && len(conns) > 0 && rollcallRequest.RequestId == rollcallResponse.RequestId {
		return false
	}

	return true
}
