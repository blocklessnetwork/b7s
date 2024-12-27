package worker

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node"
)

func (w *Worker) process(ctx context.Context, from peer.ID, msg string, payload []byte) error {

	switch msg {
	case blockless.MessageHealthCheck:
		return node.HandleMessage(ctx, from, payload, w.processHealthCheck)
	case blockless.MessageInstallFunction:
		return node.HandleMessage(ctx, from, payload, w.processInstallFunction)
	case blockless.MessageRollCall:
		return node.HandleMessage(ctx, from, payload, w.processRollCall)
	case blockless.MessageWorkOrder:
		return node.HandleMessage(ctx, from, payload, w.processWorkOrder)
	case blockless.MessageFormCluster:
		return node.HandleMessage(ctx, from, payload, w.processFormCluster)
	case blockless.MessageDisbandCluster:
		return node.HandleMessage(ctx, from, payload, w.processDisbandCluster)
	}

	return fmt.Errorf("unsupported message: %s", msg)
}
