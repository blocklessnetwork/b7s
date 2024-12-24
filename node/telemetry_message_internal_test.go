package node

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/telemetry"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
	"github.com/blocklessnetwork/b7s/testing/helpers"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_TraceHealthCheck(t *testing.T) {

	var (
		ctx = context.Background()

		peerID       = mocks.GenericPeerID
		role         = blockless.WorkerNode
		resource, _  = telemetry.CreateResource(ctx, peerID.String(), role)
		exporter, tp = helpers.CreateTracerProvider(t, resource)
		from         = mocks.GenericPeerIDs[0]

		log  = mocks.NoopLogger
		core = NewCore(log, helpers.NewLoopbackHost(t, log))
	)

	core.tracer = tracing.NewTracerFromProvider(tp, "test-tracer")

	payload := []byte(`{ "type": "MsgHealthCheck" }`)

	pre := time.Now()

	process := func(context.Context, peer.ID, string, []byte) error {
		return nil
	}

	pipeline := PubSubPipeline(blockless.DefaultTopic)
	err := core.processMessage(ctx, from, payload, pipeline, process)
	require.NoError(t, err)

	post := time.Now()

	// Now, verify the span was recorded correctly.
	spans := exporter.GetSpans()
	require.Len(t, spans, 1)

	span := spans[0]
	require.Equal(t, "MessageProcess MsgHealthCheck", span.Name)
	require.Equal(t, otelcodes.Ok, span.Status.Code)
	require.Empty(t, span.Status.Description)
	require.Equal(t, trace.SpanKindConsumer, span.SpanKind)

	// In this case we know this is a root span so verify parent field is empty.
	require.False(t, span.Parent.IsValid())

	sctx := span.SpanContext
	require.True(t, sctx.HasSpanID())
	require.True(t, sctx.HasTraceID())

	// Verify span times are correct.
	require.True(t, span.StartTime.After(pre))
	require.True(t, span.EndTime.Before(post))

	// Convert attributes to a map for easier lookup.
	attributes := attributeMap(span.Attributes)

	require.Equal(t, from.String(), attributes["message.peer"].AsString())
	require.Equal(t, pipeline.ID.String(), attributes["message.pipeline"].AsString())
	require.Equal(t, pipeline.Topic, attributes["message.topic"].AsString())
	require.Equal(t, blockless.MessageHealthCheck, attributes["message.type"].AsString())

	require.Equal(t, resource, span.Resource)
}

func attributeMap(attrs []attribute.KeyValue) map[attribute.Key]attribute.Value {

	// Convert attributes to a map for easier lookup.
	attributes := make(map[attribute.Key]attribute.Value)
	for _, attr := range attrs {
		attributes[attr.Key] = attr.Value
	}

	return attributes
}

func TestNode_ProcessedMessageMetric(t *testing.T) {

	var (
		ctx      = context.Background()
		registry = prometheus.NewRegistry()
		log      = mocks.NoopLogger

		core = NewCore(log, helpers.NewLoopbackHost(t, log))
	)

	sink, err := telemetry.CreateMetricSink(registry, telemetry.MetricsConfig{Counters: Counters})
	require.NoError(t, err)

	m, err := telemetry.CreateMetrics(sink, false)
	require.NoError(t, err)

	core.metrics = m

	// Messages to send. We will send multiple health check and disband cluster messages.
	// Note that not all messages make sense in the context of a real-world node, but we just care about having
	// a few messages flow through the system.
	var (
		// Do between 1 and 10 messages.
		limit            = 10
		healthcheckCount = rand.Intn(limit) + 1
		disbandCount     = rand.Intn(limit) + 1

		healthCheck = response.Health{}

		disbandRequest = request.DisbandCluster{
			RequestID: "request-id",
		}
	)

	msgs := []struct {
		count    int
		pipeline Pipeline
		rec      any
	}{
		{
			count:    healthcheckCount,
			pipeline: PubSubPipeline(blockless.DefaultTopic),
			rec:      healthCheck,
		},
		{
			count:    disbandCount,
			pipeline: DirectMessagePipeline,
			rec:      disbandRequest,
		},
	}

	process := func(context.Context, peer.ID, string, []byte) error {
		return nil
	}

	for _, msg := range msgs {
		for i := 0; i < msg.count; i++ {

			payload, err := json.Marshal(msg.rec)
			require.NoError(t, err)

			// We don't care if the message was processed okay (disband cluster will fail).
			_ = core.processMessage(ctx, mocks.GenericPeerID, payload, msg.pipeline, process)
		}
	}

	metricMap := helpers.MetricMap(t, registry)
	helpers.CounterCmp(t, metricMap, float64(healthcheckCount), "b7s_node_messages_processed", "type", "MsgHealthCheck")
	helpers.CounterCmp(t, metricMap, float64(disbandCount), "b7s_node_messages_processed", "type", "MsgDisbandCluster")
}
