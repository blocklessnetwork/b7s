package tracing_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
	"github.com/blocklessnetwork/b7s/testing/helpers"
)

func TestTraceInfo(t *testing.T) {

	var (
		ctx         = context.Background()
		resource, _ = telemetry.CreateResource(ctx, "instance-id", blockless.WorkerNode)
		_, tp       = helpers.CreateTracerProvider(t, resource)
		tracer      = tracing.NewTracerFromProvider(tp, fmt.Sprintf("test-tracer-%v", rand.Int()))
	)

	t.Run("untraced context returns and empty traceinfo", func(t *testing.T) {
		t.Parallel()
		ti := tracing.GetTraceInfo(context.Background())
		require.True(t, ti.Empty())
	})
	t.Run("traced context returns populated traceinfo", func(t *testing.T) {
		t.Parallel()

		childCtx, span := tracer.Start(ctx, "test-span-1")

		propagator := telemetry.CreatePropagator()
		ti := tracing.GetTraceInfoWithPropagator(childCtx, propagator)
		require.Len(t, ti.Carrier, 1)
		fields := strings.Split(ti.Carrier.Get("traceparent"), "-")
		require.Len(t, fields, 4)

		// Verify trace data.
		sctx := span.SpanContext()
		require.Equal(t, "00", fields[0]) // version field.
		require.Equal(t, sctx.TraceID().String(), fields[1])
		require.Equal(t, sctx.SpanID().String(), fields[2])
		require.Equal(t, sctx.TraceFlags().String(), fields[3])
	})
	t.Run("injected traceinfo produces identical context", func(t *testing.T) {
		t.Parallel()

		childCtx, span := tracer.Start(ctx, "test-span-2")

		propagator := telemetry.CreatePropagator()
		ti := tracing.GetTraceInfoWithPropagator(childCtx, propagator)

		newCtx := tracing.TraceContextWithPropagator(ctx, propagator, ti)
		newSpanCtx := trace.SpanContextFromContext(newCtx)

		require.Equal(t, span.SpanContext().TraceID(), newSpanCtx.TraceID())
		require.Equal(t, span.SpanContext().SpanID(), newSpanCtx.SpanID())
		require.Equal(t, span.SpanContext().TraceFlags(), newSpanCtx.TraceFlags())
	})
}
