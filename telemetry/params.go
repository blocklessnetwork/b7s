package telemetry

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	defaultResourceOpts = []resource.Option{
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("b7s"),
			semconv.ServiceVersionKey.String(vcsVersion()),
		),
	}
)

const (
	useCompressionForTraceExporters = true
	// NOTE: Temporary setting for the still young, immature stage of telemetry.
	allowInsecureTraceExporters = true
)
