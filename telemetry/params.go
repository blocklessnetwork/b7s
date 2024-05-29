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
		resource.WithContainer(), // TODO: Check if this works in docker compose cluster, and in a non-docker env with cgroups
		resource.WithAttributes(
			semconv.ServiceNameKey.String("b7s"),
			semconv.ServiceVersionKey.String(vcsVersion()),
		),
	}
)

const (
	ExporterStdout ExporterMethod = "stdout"
	ExporterGRPC   ExporterMethod = "grpc"
	ExporterHTTP   ExporterMethod = "http"
)
