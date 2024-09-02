package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

func Initialize(ctx context.Context, log zerolog.Logger, opts ...Option) (ShutdownFunc, error) {

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	err := cfg.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Setup general otel stuff.
	setupOtel(log)

	resource, err := CreateResource(ctx, cfg.ID, cfg.Role)
	if err != nil {
		return nil, fmt.Errorf("could not initialize otel resource: %w", err)
	}

	// Setup tracing.
	exporters, err := createTraceExporters(ctx, cfg.Trace)
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporters: %w", err)
	}

	if cfg.Trace.ExporterBatchTimeout == 0 {
		log.Warn().Msg("trace exporter batch timeout is disabled")
	}

	tp := CreateTracerProvider(resource, cfg.Trace.ExporterBatchTimeout, exporters...)
	otel.SetTracerProvider(tp)

	// From here on down, we have components that need shutdown.
	// Potential other shutdown functions should be appended to this slice.
	shutdownFuncs := []ShutdownFunc{
		tp.Shutdown,
	}

	// Setup metrics.
	err = initMetrics(cfg.Metrics)
	if err != nil {

		outErr := errors.Join(
			shutdownAll(shutdownFuncs)(ctx),
			fmt.Errorf("could not initialize prometheus registry: %w", err))

		return nil, outErr
	}

	shutdownFuncs = append(shutdownFuncs,
		func(context.Context) error {
			metrics.Shutdown()
			return nil
		})

	return shutdownAll(shutdownFuncs), nil
}

func shutdownAll(funcs []ShutdownFunc) ShutdownFunc {
	return func(ctx context.Context) error {
		var err error
		for _, fn := range funcs {
			err = errors.Join(err, fn(ctx))
		}
		return err
	}
}

func CreateResource(ctx context.Context, id string, role blockless.NodeRole) (*resource.Resource, error) {

	opts := append(
		defaultResourceOpts,
		resource.WithAttributes(semconv.ServiceInstanceIDKey.String(id)),
		resource.WithAttributes(b7ssemconv.ServiceRole.String(role.String())),
	)

	resource, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize resource: %w", err)
	}

	return resource, nil
}

func CreatePropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// Setup general otel stuff like logging and error handling. We will just log telemetry errors and do nothing more.
func setupOtel(log zerolog.Logger) {

	otel.SetTextMapPropagator(CreatePropagator())
	otel.SetLogger(zerologr.New(&log))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(
		func(err error) { log.Error().Err(err).Msg("telemetry error") }),
	)
}
