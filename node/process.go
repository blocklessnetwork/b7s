package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/libp2p/go-libp2p/core/peer"
	otelcodes "go.opentelemetry.io/otel/codes"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node/internal/pipeline"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

// processMessage will determine which message was received and how to process it.
func (n *Node) processMessage(ctx context.Context, from peer.ID, payload []byte, pipeline pipeline.Pipeline) (procError error) {

	// Determine message type.
	msgType, err := getMessageType(payload)
	if err != nil {
		return fmt.Errorf("could not unpack message: %w", err)
	}

	log := n.log.With().Stringer("peer", from).Str("type", msgType).Stringer("pipeline", pipeline).Logger()

	if !messageAllowedOnPipeline(msgType, pipeline) {
		log.Debug().Msg("message not allowed on pipeline")
		return nil
	}

	if !n.messageAllowedForRole(msgType) {
		log.Debug().Msg("message not intended for our role")
		return nil
	}

	n.metrics.IncrCounterWithLabels(messagesProcessedMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
	defer func() {
		switch procError {
		case nil:
			n.metrics.IncrCounterWithLabels(messagesProcessedOkMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
		default:
			n.metrics.IncrCounterWithLabels(messagesProcessedErrMetric, 1, []metrics.Label{{Name: "type", Value: msgType}})
		}
	}()

	ctx, err = tracing.TraceContextFromMessage(ctx, payload)
	if err != nil {
		n.log.Error().Err(err).Msg("could not get trace context from message")
	}

	ctx, span := n.tracer.Start(ctx, msgProcessSpanName(msgType), msgProcessSpanOpts(from, msgType, pipeline)...)
	defer span.End()
	// NOTE: This function checks the named return error value in order to set the span status accordingly.
	defer func() {
		if procError == nil {
			span.SetStatus(otelcodes.Ok, spanStatusOK)
			return
		}

		if allowErrorLeakToTelemetry {
			span.SetStatus(otelcodes.Error, procError.Error())
			return
		}

		span.SetStatus(otelcodes.Error, spanStatusErr)
	}()

	log.Debug().Msg("received message from peer")

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

func (n *Node) messageAllowedForRole(msgType string) bool {

	// Worker node allowed messages.
	if n.isWorker() {
		switch msgType {
		case blockless.MessageHealthCheck,
			blockless.MessageInstallFunction,
			blockless.MessageRollCall,
			blockless.MessageExecute,
			blockless.MessageFormCluster,
			blockless.MessageDisbandCluster:
			return true

		default:
			return false
		}
	}

	// Head node allowed messages.
	switch msgType {

	case blockless.MessageHealthCheck,
		blockless.MessageInstallFunctionResponse,
		blockless.MessageRollCallResponse,
		blockless.MessageExecute,
		blockless.MessageExecuteResponse,
		blockless.MessageFormClusterResponse:

		// NOTE: We provide a mechanism via the REST API to broadcast function install, so there's a case for this being supported.
		return true

	default:
		return false
	}
}

func handleMessage[T blockless.Message](ctx context.Context, from peer.ID, payload []byte, processFunc func(ctx context.Context, from peer.ID, msg T) error) error {

	var msg T
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return fmt.Errorf("could not unmarshal message: %w", err)
	}

	// If the message provides a validation mechanism - use it.
	type validator interface {
		Valid() error
	}

	vmsg, ok := any(msg).(validator)
	if ok {
		err = vmsg.Valid()
		if err != nil {
			return fmt.Errorf("rejecting message that failed validation: %w", err)
		}
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
