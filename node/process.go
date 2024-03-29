package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

// TODO: Set up a chain: message ID => model => handler

// processMessage will determine which message was received and how to process it.
func (n *Node) processMessage(ctx context.Context, from peer.ID, payload []byte) error {

	// Determine message type.
	msg, err := unpackMessage(payload)
	if err != nil {
		return fmt.Errorf("could not unpack message: %w", err)
	}

	msgType := msg.Type()

	n.log.Trace().Str("peer", from.String()).Str("message", msgType).Msg("received message from peer")

	// Get the registered handler for the message.
	switch msgType {

	case blockless.MessageHealthCheck:
		return n.processHealthCheck(ctx, from, msg.(response.Health))

	case blockless.MessageInstallFunction:
		return n.processInstallFunction(ctx, from, msg.(request.InstallFunction))
	case blockless.MessageInstallFunctionResponse:
		return n.processInstallFunctionResponse(ctx, from, msg.(response.InstallFunction))

	case blockless.MessageRollCall:
		return n.processRollCall(ctx, from, msg.(request.RollCall))
	case blockless.MessageRollCallResponse:
		return n.processRollCallResponse(ctx, from, msg.(response.RollCall))

	case blockless.MessageExecute:
		return n.processExecute(ctx, from, msg.(request.Execute))
	case blockless.MessageExecuteResponse:
		return n.processExecuteResponse(ctx, from, msg.(response.Execute))

	case blockless.MessageFormCluster:
		return n.processFormCluster(ctx, from, msg.(request.FormCluster))
	case blockless.MessageFormClusterResponse:
		return n.processFormClusterResponse(ctx, from, msg.(response.FormCluster))
	case blockless.MessageDisbandCluster:
		return n.processDisbandCluster(ctx, from, msg.(request.DisbandCluster))

	default:
		return fmt.Errorf("unsupported message type (from: %s): %s", from.String(), msgType)
	}
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

func unpackMessage(payload []byte) (blockless.Message, error) {

	// Determine message type.
	msgType, err := getMessageType(payload)
	if err != nil {
		return nil, fmt.Errorf("could not determine message type: %w", err)
	}

	switch msgType {

	case blockless.MessageHealthCheck:
		return unmarshalJSON[response.Health](payload)
	case blockless.MessageInstallFunction:
		return unmarshalJSON[request.InstallFunction](payload)
	case blockless.MessageInstallFunctionResponse:
		return unmarshalJSON[response.InstallFunction](payload)
	case blockless.MessageRollCall:
		return unmarshalJSON[request.RollCall](payload)
	case blockless.MessageRollCallResponse:
		return unmarshalJSON[response.RollCall](payload)
	case blockless.MessageExecute:
		return unmarshalJSON[request.Execute](payload)
	case blockless.MessageExecuteResponse:
		return unmarshalJSON[response.Execute](payload)
	case blockless.MessageFormCluster:
		return unmarshalJSON[request.FormCluster](payload)
	case blockless.MessageFormClusterResponse:
		return unmarshalJSON[response.FormCluster](payload)
	case blockless.MessageDisbandCluster:
		return unmarshalJSON[request.DisbandCluster](payload)

	default:
		return nil, fmt.Errorf("unknown message type: %w", err)
	}
}

func unmarshalJSON[T any](payload []byte) (T, error) {

	var obj T
	err := json.Unmarshal(payload, &obj)
	if err != nil {
		return obj, fmt.Errorf("could not unmarshal message: %w", err)
	}

	return obj, nil
}
