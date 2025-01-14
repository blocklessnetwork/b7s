package telemetry

import (
	"github.com/blessnetwork/b7s/info"
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
			semconv.ServiceVersionKey.String(info.VcsVersion()),
		),
	}
)

const (
	useCompressionForTraceExporters = true
	// NOTE: Temporary setting for the still young, immature stage of telemetry.
	allowInsecureTraceExporters = true

	metricPrefix = "b7s"
)
