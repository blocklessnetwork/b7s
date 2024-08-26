package tracing_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
	"github.com/blocklessnetwork/b7s/testing/helpers"
)

func TestTracer_TraceFunction(t *testing.T) {

	var (
		tracerName = "test-tracer"
		fnErr      = errors.New("function-error")
	)

	resource, err := telemetry.CreateResource(context.Background(), "instance-id", blockless.WorkerNode)
	require.NoError(t, err)

	tests := []struct {
		name      string
		wantedErr error
	}{
		{
			name:      "function failure returns span with error",
			wantedErr: fnErr,
		},
		{
			name:      "function success returns span ok",
			wantedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			exporter, tp := helpers.CreateTracerProvider(t, resource)

			var (
				spanName = fmt.Sprintf("test-span-%v", rand.Int())

				keys = []string{
					fmt.Sprintf("attr-key-1-%v", rand.Int()),
					fmt.Sprintf("attr-key-2-%v", rand.Int()),
				}
				values = []any{
					fmt.Sprintf("attr-value-1-%v", rand.Int()),
					rand.Int(),
				}

				attrs = []attribute.KeyValue{
					{
						Key:   attribute.Key(keys[0]),
						Value: attribute.StringValue(values[0].(string)),
					},
					{
						Key:   attribute.Key(keys[1]),
						Value: attribute.IntValue(values[1].(int)),
					},
				}
			)

			// Function that will be executed and traced.
			executed := false
			fn := func() error {
				executed = true
				return test.wantedErr
			}

			err := tracing.NewTracerFromProvider(tp, tracerName).WithSpanFromContext(context.Background(), spanName, fn, trace.WithAttributes(attrs...))
			require.ErrorIs(t, err, test.wantedErr)

			require.True(t, executed)

			spans := exporter.GetSpans()
			require.Len(t, spans, 1)

			span := spans[0]
			require.Equal(t, spanName, span.Name)

			switch test.wantedErr {
			case nil:
				require.Equal(t, span.Status.Code, otelcodes.Ok)
			default:
				require.Equal(t, span.Status.Code, otelcodes.Error)
				require.Equal(t, fnErr.Error(), span.Status.Description)
			}

			require.Len(t, span.Attributes, 2)
			require.Equal(t, string(span.Attributes[0].Key), keys[0])
			require.Equal(t, values[0].(string), span.Attributes[0].Value.AsString())
			require.Equal(t, string(span.Attributes[1].Key), keys[1])
			require.Equal(t, int64(values[1].(int)), span.Attributes[1].Value.AsInt64())
			require.Equal(t, tracerName, span.InstrumentationLibrary.Name)
		})
	}
}
