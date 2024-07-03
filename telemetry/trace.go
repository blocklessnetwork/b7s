package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

func newTracerProvider(ctx context.Context, cfg Config) (*trace.TracerProvider, error) {

	// TODO: Fix hardcoded true.
	exporter, err := traceExporter(ctx, cfg.ExporterMethod, true)
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporter: %w", err)
	}

	// TODO: Does this correctly join the ID with the remaining stuff?
	opts := defaultResourceOpts
	if cfg.ID != "" {
		opts = append(opts, resource.WithAttributes(semconv.ServiceInstanceIDKey.String(cfg.ID)))
	}
	if cfg.Role.Valid() {
		opts = append(opts, resource.WithAttributes(b7ssemconv.ServiceRole.String(cfg.Role.String())))
	}

	resource, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize resource: %w", err)
	}

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(cfg.BatchTraceTimeout)),
		trace.WithResource(resource),
	)

	return provider, nil
}

func newPropagator() propagation.TextMapPropagator {

	pp := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return pp
}

// TODO: What do we need in the config?
// GRPC:
// - endpoint
// - use compression by default
// - TLS credentials
// Perhaps:
// - insecure?
//
// HTTP:
// - endpoint
// - compression (use by default)
// - TLS credentials
// Perhaps:
// - insecure
func traceExporter(ctx context.Context, m ExporterMethod, allowInsecure bool) (trace.SpanExporter, error) {

	switch m {
	case ExporterGRPC:
		opts := []otlptracegrpc.Option{otlptracegrpc.WithCompressor("gzip")}
		if allowInsecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		return otlptracegrpc.New(ctx, opts...)

	case ExporterHTTP:

		opts := []otlptracehttp.Option{otlptracehttp.WithCompression(otlptracehttp.GzipCompression)}
		if allowInsecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		return otlptracehttp.New(ctx, opts...)

	// NOTE: STDOUT exporterr is not for production use.
	case ExporterStdout:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	default:
		return nil, errors.New("unsupported exporter type")
	}
}
