package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/zerologr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

var (
	globalRegistry *prometheus.Registry
)

func Initialize(ctx context.Context, log zerolog.Logger, opts ...Option) (shutdown ShutdownFunc, err error) {

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	shutdown, err = initializeTracing(ctx, log, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not initialize tracing: %w", err)
	}

	err = initPrometheusRegistry()
	if err != nil {
		shutdown(ctx)
		return nil, fmt.Errorf("could not initialize prometheus registry: %w", err)
	}

	return shutdown, nil
}

func initializeTracing(ctx context.Context, log zerolog.Logger, cfg Config) (shutdown ShutdownFunc, err error) {

	var shutdownFuncs []ShutdownFunc
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set logger and global error handler function - just log error and nothing else.
	otel.SetErrorHandler(otel.ErrorHandlerFunc(
		func(err error) { log.Error().Err(err).Msg("telemetry error") },
	))
	otel.SetLogger(zerologr.New(&log))

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tp, err := newTracerProvider(ctx, cfg)
	if err != nil {
		handleErr(fmt.Errorf("could not create new trace provider: %w", err))
		return
	}
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
	otel.SetTracerProvider(tp)

	return
}
