package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

func newTracerProvider(ctx context.Context, cfg Config) (*trace.TracerProvider, error) {

	// Setup resource.
	opts := defaultResourceOpts
	if cfg.ID == "" {
		return nil, errors.New("instance ID is required")
	}
	opts = append(opts, resource.WithAttributes(semconv.ServiceInstanceIDKey.String(cfg.ID)))

	if !cfg.Role.Valid() {
		return nil, errors.New("node role is required")
	}
	opts = append(opts, resource.WithAttributes(b7ssemconv.ServiceRole.String(cfg.Role.String())))

	resource, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize resource: %w", err)
	}

	// Setup exporters.
	exporters, err := traceExporters(ctx, cfg.Trace)
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporter: %w", err)
	}

	traceOpts := []trace.TracerProviderOption{trace.WithResource(resource)}
	for _, exporter := range exporters {
		traceOpts = append(traceOpts, trace.WithBatcher(exporter, trace.WithBatchTimeout(cfg.Trace.ExporterBatchTimeout)))
	}

	return trace.NewTracerProvider(traceOpts...), nil
}

func newPropagator() propagation.TextMapPropagator {

	pp := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	return pp
}

func traceExporters(ctx context.Context, tcfg TraceConfig) ([]trace.SpanExporter, error) {

	var exporters []trace.SpanExporter
	if tcfg.GRPC.Enabled {

		ex, err := newGRPCExporter(ctx, tcfg.GRPC)
		if err != nil {
			return nil, fmt.Errorf("could not create new GRPC exporter: %w", err)
		}

		exporters = append(exporters, ex)
	}

	if tcfg.HTTP.Enabled {

		ex, err := newHTTPExporter(ctx, tcfg.HTTP)
		if err != nil {
			return nil, fmt.Errorf("could not create new GRPC exporter: %w", err)
		}

		exporters = append(exporters, ex)
	}

	return exporters, nil
}

func newGRPCExporter(ctx context.Context, cfg TraceGRPCConfig) (*otlptrace.Exporter, error) {

	var opts []otlptracegrpc.Option
	if cfg.UseCompression {
		opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	}

	if cfg.AllowInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	return otlptracegrpc.New(ctx, opts...)
}

func newHTTPExporter(ctx context.Context, cfg TraceHTTPConfig) (*otlptrace.Exporter, error) {

	var opts []otlptracehttp.Option
	if cfg.UseCompression {
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression)
	}

	if cfg.AllowInsecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	return otlptracehttp.New(ctx, opts...)
}
