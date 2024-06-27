package tracing

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TraceInfo struct {
	Carrier propagation.MapCarrier
}

func GetTraceInfo(ctx context.Context) TraceInfo {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	return TraceInfo{Carrier: carrier}
}

func TraceContextFromMessage(ctx context.Context, payload []byte) (context.Context, error) {

	var traceInfo TraceInfo
	err := json.Unmarshal(payload, &traceInfo)
	if err != nil {
		return ctx, fmt.Errorf("could not extract trace info from context: %w", err)
	}

	return TraceContext(ctx, traceInfo), nil
}

func TraceContext(ctx context.Context, t TraceInfo) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, t.Carrier)
}
