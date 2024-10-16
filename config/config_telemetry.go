package config

import (
	"time"
)

type Telemetry struct {
	Tracing Tracing `koanf:"tracing"`
	Metrics Metrics `koanf:"metrics"`
}

type Tracing struct {
	Enable               bool          `koanf:"enable" flag:"enable-tracing"`
	ExporterBatchTimeout time.Duration `koanf:"exporter-batch-timeout"`
	GRPC                 GRPCTracing   `koanf:"grpc"`
	HTTP                 HTTPTracing   `koanf:"http"`
}

type GRPCTracing struct {
	Endpoint string `koanf:"endpoint" flag:"tracing-grpc-endpoint"`
}

type HTTPTracing struct {
	Endpoint string `koanf:"endpoint" flag:"tracing-http-endpoint"`
}

type Metrics struct {
	Enable            bool   `koanf:"enable" flag:"enable-metrics"`
	PrometheusAddress string `koanf:"prometheus-address" flag:"prometheus-address"`
}
