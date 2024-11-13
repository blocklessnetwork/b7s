package node

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/telemetry"
	"github.com/blocklessnetwork/b7s/testing/helpers"
)

func TestTelemetry_SaveTraceContext(t *testing.T) {

	var (
		ctx         = context.Background()
		resource, _ = telemetry.CreateResource(ctx, "instance-id", blockless.WorkerNode)
		_, tp       = helpers.CreateTracerProvider(t, resource)
		tracer      = tp.Tracer("test-tracer")
		spanName    = fmt.Sprintf("test-span-%v", rand.Int())
	)

	// TODO: Not too pretty to have this in tests.
	otel.SetTextMapPropagator(telemetry.CreatePropagator())

	childCtx, span := tracer.Start(ctx, spanName)
	sctx := span.SpanContext()

	msg := new(response.Execute)
	saveTraceContext(childCtx, msg)

	fields := strings.Split(msg.BaseMessage.TraceInfo.Carrier["traceparent"], "-")

	require.Len(t, fields, 4)
	require.Equal(t, "00", fields[0])
	require.Equal(t, sctx.TraceID().String(), fields[1])
	require.Equal(t, sctx.SpanID().String(), fields[2])
	require.Equal(t, sctx.TraceFlags().String(), fields[3])
}
