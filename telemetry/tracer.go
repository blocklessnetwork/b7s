package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer trace.Tracer
}

func NewTracer(name string) *Tracer {

	return &Tracer{
		tracer: otel.Tracer(name),
	}
}

func (t *Tracer) WithSpanFromContext(ctx context.Context, spanName string, f func() error, opts ...trace.SpanStartOption) error {

	_, span := t.tracer.Start(ctx, spanName, opts...)
	defer span.End()

	err := f()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}
