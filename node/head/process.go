package head

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node"
)

// processMessage will determine which message was received and how to process it.
func (h *HeadNode) process(ctx context.Context, from peer.ID, msg string, payload []byte) error {

	switch msg {
	case blockless.MessageHealthCheck:
		return node.HandleMessage(ctx, from, payload, h.processHealthCheck)
	case blockless.MessageInstallFunctionResponse:
		return node.HandleMessage(ctx, from, payload, h.processInstallFunctionResponse)
	case blockless.MessageExecute:
		return node.HandleMessage(ctx, from, payload, h.processExecute)
	case blockless.MessageRollCallResponse:
		return node.HandleMessage(ctx, from, payload, h.processRollCallResponse)
	case blockless.MessageWorkOrderResponse:
		return node.HandleMessage(ctx, from, payload, h.processWorkOrderResponse)
	case blockless.MessageFormClusterResponse:
		return node.HandleMessage(ctx, from, payload, h.processFormClusterResponse)
	}

	return fmt.Errorf("unsupported message: %s", msg)
}
