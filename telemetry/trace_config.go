package telemetry

import (
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

var DefaultTraceConfig = TraceConfig{
	ExporterBatchTimeout: 1 * time.Second,
	GRPC: TraceGRPCConfig{
		Enabled:        false,
		AllowInsecure:  allowInsecureTraceExporters,
		UseCompression: useCompressionForTraceExporters,
	},
	HTTP: TraceHTTPConfig{
		Enabled:        false,
		AllowInsecure:  allowInsecureTraceExporters,
		UseCompression: useCompressionForTraceExporters,
	},
}

// TODO: Update trace exporters configs
// GRPC, HTTP:
// - TLS credentials
// - disable insecure when mature
type TraceConfig struct {
	// Node ID, registered as service instance ID attribute.
	ID string
	// Node role, registered as service role attribute.
	Role blockless.NodeRole
	// Maximum time after which exporters will send batched span.
	ExporterBatchTimeout time.Duration
	// Configuration for GRPC trace exporter.
	GRPC TraceGRPCConfig
	// Configuration for HTTP trace exporter.
	HTTP TraceHTTPConfig
	// Configuration for the InMem trace exporter (used for testing mainly).
	InMem TraceInMemConfig
}

type TraceGRPCConfig struct {
	Enabled        bool
	Endpoint       string
	AllowInsecure  bool
	UseCompression bool
	// TLSConfig
}

type TraceHTTPConfig struct {
	Enabled        bool
	Endpoint       string
	AllowInsecure  bool
	UseCompression bool
	// TLSConfig
}

type TraceInMemConfig struct {
	Enabled bool
}

type TraceOption func(*TraceConfig)

func WithNodeRole(r blockless.NodeRole) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.Role = r
	}
}

func WithBatchTraceTimeout(t time.Duration) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.ExporterBatchTimeout = t
	}
}

func WithID(id string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.ID = id
	}
}

func WithGRPCTracing(endpoint string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.GRPC.Endpoint = endpoint
		cfg.GRPC.Enabled = endpoint != ""
	}
}

func WithHTTPTracing(endpoint string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.HTTP.Endpoint = endpoint
		cfg.HTTP.Enabled = endpoint != ""
	}
}
