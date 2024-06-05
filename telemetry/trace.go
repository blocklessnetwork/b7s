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

	exporter, err := traceExporter(ctx, cfg.ExporterMethod)
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

func traceExporter(ctx context.Context, m ExporterMethod) (trace.SpanExporter, error) {

	// TODO: Allow insecure - from config.
	switch m {
	case ExporterGRPC:
		return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	case ExporterStdout:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	case ExporterHTTP:
		return otlptracehttp.New(ctx, otlptracehttp.WithInsecure())

	default:
		return nil, errors.New("unsupported exporter type")
	}
}
