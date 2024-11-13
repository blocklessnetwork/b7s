package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// Create a new tracer provider.
// NOTE: batchTimeout should not be set to zero for production use.
func CreateTracerProvider(resource *resource.Resource, batchTimeout time.Duration, exporters ...trace.SpanExporter) *trace.TracerProvider {

	opts := []trace.TracerProviderOption{
		trace.WithResource(resource),
	}

	for _, exporter := range exporters {
		if batchTimeout == 0 {
			opts = append(opts, trace.WithSyncer(exporter))
			continue
		}

		opts = append(opts, trace.WithBatcher(exporter, trace.WithBatchTimeout(batchTimeout)))
	}

	return trace.NewTracerProvider(opts...)
}

func createTraceExporters(ctx context.Context, tcfg TraceConfig) ([]trace.SpanExporter, error) {

	var exporters []trace.SpanExporter

	// If creating some of the exporters fails, shutdown others that were created.
	shutdown := func() {
		for _, ex := range exporters {
			_ = ex.Shutdown(ctx)
		}
	}

	if tcfg.GRPC.Enabled {

		ex, err := NewGRPCExporter(ctx, tcfg.GRPC)
		if err != nil {
			return nil, fmt.Errorf("could not create new GRPC exporter: %w", err)
		}

		exporters = append(exporters, ex)
	}

	if tcfg.HTTP.Enabled {

		ex, err := NewHTTPExporter(ctx, tcfg.HTTP)
		if err != nil {
			shutdown()
			return nil, fmt.Errorf("could not create new HTTP exporter: %w", err)
		}

		exporters = append(exporters, ex)
	}

	if tcfg.InMem.Enabled {
		exporters = append(exporters, NewInMemExporter())
	}

	return exporters, nil
}

func NewGRPCExporter(ctx context.Context, cfg TraceGRPCConfig) (*otlptrace.Exporter, error) {

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.UseCompression {
		opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	}

	if cfg.AllowInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	return otlptracegrpc.New(ctx, opts...)
}

func NewHTTPExporter(ctx context.Context, cfg TraceHTTPConfig) (*otlptrace.Exporter, error) {

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Endpoint),
	}
	if cfg.UseCompression {
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression)
	}

	if cfg.AllowInsecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	return otlptracehttp.New(ctx, opts...)
}

func NewInMemExporter() *tracetest.InMemoryExporter {
	return tracetest.NewInMemoryExporter()
}
