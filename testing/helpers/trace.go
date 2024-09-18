package helpers

import (
	"testing"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/blocklessnetwork/b7s/telemetry"
)

func CreateTracerProvider(t *testing.T, resource *resource.Resource) (*tracetest.InMemoryExporter, *trace.TracerProvider) {
	t.Helper()

	exporter := telemetry.NewInMemExporter()
	tp := telemetry.CreateTracerProvider(resource, 0, exporter)

	return exporter, tp
}
