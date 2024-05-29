package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
)

func SetupSDK(ctx context.Context, opts ...Option) (shutdown ShutdownFunc, err error) {

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

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTracerProvider(ctx, cfg)
	if err != nil {
		handleErr(fmt.Errorf("could not create new trace provider: %w", err))
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// TODO: meter provider
	// TODO: logger provider
	return
}
