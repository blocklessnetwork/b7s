package node

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
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

		node = createNode(t, role)
	)

	node.tracer = tracing.NewTracerFromProvider(tp, "test-tracer")

	payload := []byte(`{ "type": "MsgHealthCheck" }`)

	pre := time.Now()

	err := node.processMessage(ctx, from, payload, subscriptionPipeline)
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
	require.Equal(t, subscriptionPipeline.String(), attributes["message.pipeline"].AsString())
	require.Equal(t, blockless.MessageHealthCheck, attributes["message.type"].AsString())

	require.Equal(t, resource, span.Resource)
}

func TestNode_TraceExecution(t *testing.T) {

	// This is a more involved test, somewhere close to an integration test. It covers the scenario of a worker node processing an execution request,
	// where it creates a span for the execution request, invokes the executor, and sends the response to the client. We create a trace exporter to
	// collect these traces and verify their correctness and completness (span attributes-wise), as well as their relationship - executor and
	// message sending spans should be children of the original span created for the execution.

	var (
		ctx = context.Background()

		peerID       = mocks.GenericPeerID
		role         = blockless.WorkerNode
		resource, _  = telemetry.CreateResource(ctx, peerID.String(), role)
		exporter, tp = helpers.CreateTracerProvider(t, resource)

		executorTestSpan = fmt.Sprintf("test-span-%v", rand.Int())
	)

	// Create worker node.
	node := createNode(t, role)
	tracer := tracing.NewTracerFromProvider(tp, "test-tracer")
	node.tracer = tracer

	// Prepare a side to receive the response.
	receiver, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)
	hostAddNewPeer(t, node.host, receiver)

	// Register a processor just so we receive the message.
	receiver.SetStreamHandler(blockless.ProtocolID, func(_ network.Stream) {})

	tests := []struct {
		name    string
		req     request.Execute
		wantErr error
	}{
		{
			name: "execution failed",
			req: request.Execute{
				RequestID: newRequestID(),
				Request: execute.Request{
					FunctionID: fmt.Sprintf("test-function-cid-1-%v", rand.Int()),
					Method:     fmt.Sprintf("test-method-%v.wasm", rand.Int()),
					Config: execute.Config{
						NodeCount:          rand.Intn(16),
						ConsensusAlgorithm: "",
					},
				},
			},
			wantErr: errors.New("test-error"),
		},
		{
			name: "execution ok",
			req: request.Execute{
				RequestID: newRequestID(),
				Request: execute.Request{
					FunctionID: fmt.Sprintf("test-function-cid-2-%v", rand.Int()),
					Method:     fmt.Sprintf("test-method-%v.wasm", rand.Int()),
					Config: execute.Config{
						NodeCount:          rand.Intn(16),
						ConsensusAlgorithm: "",
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			executor := mocks.BaselineExecutor(t)
			// Populate the span in the executor and verify span chain (parent).
			executor.ExecFunctionFunc = func(ctx context.Context, requestID string, req execute.Request) (execute.Result, error) {
				_, span := tracer.Start(ctx, executorTestSpan)
				defer span.End()

				return execute.Result{}, test.wantErr
			}

			node.executor = executor

			verifySendSpan := func(t *testing.T, span tracetest.SpanStub) {
				t.Helper()

				require.True(t, span.Parent.IsValid())

				attributes := attributeMap(span.Attributes)

				require.Equal(t, receiver.ID().String(), attributes["message.peer"].AsString())
				require.Equal(t, directMessagePipeline.String(), attributes["message.pipeline"].AsString())

			}
			verifyExecutionSpan := func(t *testing.T, span tracetest.SpanStub) {
				t.Helper()

				// In this case we know this is a root span so verify parent field is empty.
				require.False(t, span.Parent.IsValid())

				attributes := attributeMap(span.Attributes)

				require.Equal(t, test.req.Request.FunctionID, attributes["function.cid"].AsString())
				require.Equal(t, test.req.Request.Method, attributes["function.method"].AsString())
				require.Equal(t, int64(test.req.Request.Config.NodeCount), attributes["execution.node.count"].AsInt64())
				require.Equal(t, test.req.Request.Config.ConsensusAlgorithm, attributes["execution.consensus"].AsString())
				require.Equal(t, test.req.RequestID, attributes["execution.request.id"].AsString())

				require.Equal(t, resource, span.Resource)
			}

			pre := time.Now()

			err := node.workerProcessExecute(ctx, receiver.ID(), test.req)
			require.NoError(t, err)

			post := time.Now()

			// Now, verify the spans were recorded correctly.
			// We expect at minimum spans for:
			// 1. workerProcessExecute function
			// 2. executor
			// 3. message sending

			spans := make(map[string]tracetest.SpanStub)
			for _, span := range exporter.GetSpans() {

				// Verify all span times are correct.
				require.Truef(t, span.StartTime.After(pre), "unexpected timestamp, pre: %s, span start time: %s", pre.Format(time.RFC3339), span.StartTime.Format(time.RFC3339))
				require.Truef(t, span.EndTime.Before(post), "unexpected timestamp, post: %s, span end time: %s", post.Format(time.RFC3339), span.EndTime.Format(time.RFC3339))

				// Verify span IDs are valid.
				sctx := span.SpanContext
				require.True(t, sctx.HasSpanID())
				require.True(t, sctx.HasTraceID())

				spans[span.Name] = span
			}

			// Verify worker execute span.
			workerExecuteSpan, ok := spans[spanWorkerExecute]
			require.True(t, ok)
			verifyExecutionSpan(t, workerExecuteSpan)

			// Verify executor span.
			span, ok := spans[executorTestSpan]
			require.True(t, ok)

			require.True(t, span.Parent.IsValid())
			require.True(t, span.Parent.Equal(workerExecuteSpan.SpanContext))

			// Verify send span.
			span, ok = spans[msgSendSpanName(spanMessageSend, blockless.MessageExecuteResponse)]
			require.True(t, ok)
			verifySendSpan(t, span)

			// Clear state of the exporter so subsequent tests have a clear slate.
			exporter.Reset()
		})
	}
}

func attributeMap(attrs []attribute.KeyValue) map[attribute.Key]attribute.Value {

	// Convert attributes to a map for easier lookup.
	attributes := make(map[attribute.Key]attribute.Value)
	for _, attr := range attrs {
		attributes[attr.Key] = attr.Value
	}

	return attributes
}
