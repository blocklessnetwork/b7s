package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// processMessage will determine which message was received and how to process it.
func (n *Node) processMessage(ctx context.Context, from peer.ID, payload []byte) error {

	// Determine message type.
	msgType, err := getMessageType(payload)
	if err != nil {
		return fmt.Errorf("could not unpack message: %w", err)
	}

	n.log.Trace().Str("peer", from.String()).Str("type", msgType).Msg("received message from peer")

	switch msgType {
	case blockless.MessageHealthCheck:
		return handleMessage(ctx, from, payload, n.processHealthCheck)

	case blockless.MessageInstallFunction:
		return handleMessage(ctx, from, payload, n.processInstallFunction)
	case blockless.MessageInstallFunctionResponse:
		return handleMessage(ctx, from, payload, n.processInstallFunctionResponse)

	case blockless.MessageRollCall:
		return handleMessage(ctx, from, payload, n.processRollCall)
	case blockless.MessageRollCallResponse:
		return handleMessage(ctx, from, payload, n.processRollCallResponse)

	case blockless.MessageExecute:
		return handleMessage(ctx, from, payload, n.processExecute)
	case blockless.MessageExecuteResponse:
		return handleMessage(ctx, from, payload, n.processExecuteResponse)

	case blockless.MessageFormCluster:
		return handleMessage(ctx, from, payload, n.processFormCluster)
	case blockless.MessageFormClusterResponse:
		return handleMessage(ctx, from, payload, n.processFormClusterResponse)
	case blockless.MessageDisbandCluster:
		return handleMessage(ctx, from, payload, n.processDisbandCluster)

	default:
		return fmt.Errorf("unknown message type: %s", msgType)
	}
}

func handleMessage[T blockless.Message](ctx context.Context, from peer.ID, payload []byte, processFunc func(ctx context.Context, from peer.ID, msg T) error) error {

	var msg T
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return fmt.Errorf("could not unmarshal message: %w", err)
	}

	return processFunc(ctx, from, msg)
}

// getMessageType will return the `type` string field from the JSON payload.
func getMessageType(payload []byte) (string, error) {

	type baseMessage struct {
		Type string `json:"type,omitempty"`
	}
	var message baseMessage
	err := json.Unmarshal(payload, &message)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal message: %w", err)
	}

	return message.Type, nil
}
