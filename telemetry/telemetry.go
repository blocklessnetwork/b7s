package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

func SetupSDK(ctx context.Context, log zerolog.Logger, opts ...Option) (shutdown ShutdownFunc, err error) {

	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

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

	// Set global error handler function - just log error and nothing else.
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

	// TODO: meter provider
	// TODO: logger provider
	return
}
