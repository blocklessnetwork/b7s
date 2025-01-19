package tracing_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/telemetry"
	"github.com/blessnetwork/b7s/telemetry/tracing"
	"github.com/blessnetwork/b7s/testing/helpers"
)

func TestTracer_TraceFunction(t *testing.T) {

	var (
		tracerName = "test-tracer"
		fnErr      = errors.New("function-error")
	)

	resource, err := telemetry.CreateResource(context.Background(), "instance-id", bls.WorkerNode)
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
				attrs    = createAttributes()
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

			require.Equal(t, span.Attributes, attrs)

			switch test.wantedErr {
			case nil:
				require.Equal(t, span.Status.Code, otelcodes.Ok)
			default:
				require.Equal(t, span.Status.Code, otelcodes.Error)
				require.Equal(t, fnErr.Error(), span.Status.Description)
			}
		})
	}
}

func createAttributes() []attribute.KeyValue {

	keys := []string{
		fmt.Sprintf("test-attr-key-1-%v", rand.Int()),
		fmt.Sprintf("test-attr-key-2-%v", rand.Int()),
		fmt.Sprintf("test-attr-key-3-%v", rand.Int()),
	}

	values := []any{
		fmt.Sprintf("attr-value-1-%v", rand.Int()),
		rand.Int(),
		rand.Float64(),
	}

	attrs := []attribute.KeyValue{
		{
			Key:   attribute.Key(keys[0]),
			Value: attribute.StringValue(values[0].(string)),
		},
		{
			Key:   attribute.Key(keys[1]),
			Value: attribute.IntValue(values[1].(int)),
		},
		{
			Key:   attribute.Key(keys[2]),
			Value: attribute.Float64Value(values[2].(float64)),
		},
	}

	return attrs
}
