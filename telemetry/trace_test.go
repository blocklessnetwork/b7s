package telemetry_test

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
)

func TestTelemetry_TraceProviderInMem(t *testing.T) {

	var (
		ctx = context.Background()

		spanName   = fmt.Sprintf("test-span-%v", rand.Int())
		attrValue  = fmt.Sprintf("span-attr-%v", rand.Int())
		spanError  = errors.New("test-error")
		resourceID = "instance-id"
		role       = bls.WorkerNode
	)

	resource, err := telemetry.CreateResource(ctx, resourceID, role)
	require.NoError(t, err)

	exporter := telemetry.NewInMemExporter()
	tp := telemetry.CreateTracerProvider(resource, 0, exporter)
	defer tp.Shutdown(ctx)

	tracer := tp.Tracer("test")

	traceFunc := func() (retErr error) {
		_, span := tracer.Start(ctx, spanName, trace.WithAttributes(attribute.Key("span-key").String(attrValue)))
		defer span.End()

		defer func() {
			switch retErr {
			case nil:
				span.SetStatus(otelcodes.Ok, "")
			default:
				span.SetStatus(otelcodes.Error, retErr.Error())
			}
		}()

		return spanError
	}

	traceFunc()

	spans := exporter.GetSpans()
	require.NotEmpty(t, spans)

	found := false
	for _, span := range spans {
		if span.Name != spanName {
			continue
		}

		found = true

		// Already verified span name is correct.

		require.Equal(t, span.Resource, resource)

		require.Equal(t, span.Status.Code, otelcodes.Error)
		require.Equal(t, spanError.Error(), span.Status.Description)

		require.Len(t, span.Attributes, 1)
		require.Equal(t, "span-key", string(span.Attributes[0].Key))
		require.Equal(t, attrValue, span.Attributes[0].Value.AsString())

	}

	require.True(t, found)
}
