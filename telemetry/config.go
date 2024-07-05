package telemetry

import (
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

var DefaultConfig = Config{

	Trace: TraceConfig{
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
	},
}

type Config struct {
	// Node ID, registered as service instance ID attribute.
	ID string
	// Node role, registered as service role attribute.
	Role blockless.NodeRole
	// Tracer configuration.
	Trace TraceConfig
}

// TODO: Update trace exporters configs
// GRPC, HTTP:
// - TLS credentials
// - disable insecure when mature
type TraceConfig struct {
	// Maximum time after which exporters will send batched span.
	ExporterBatchTimeout time.Duration
	// Configuration for GRPC trace exporter.
	GRPC TraceGRPCConfig
	// Configuration for HTTP trace exporter.
	HTTP TraceHTTPConfig
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

type Option func(*Config)

func WithNodeRole(r blockless.NodeRole) Option {
	return func(cfg *Config) {
		cfg.Role = r
	}
}

func WithBatchTraceTimeout(t time.Duration) Option {
	return func(cfg *Config) {
		cfg.Trace.ExporterBatchTimeout = t
	}
}

func WithID(id string) Option {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

func WithGRPCTracing(endpoint string) Option {
	return func(cfg *Config) {
		cfg.Trace.GRPC.Endpoint = endpoint
		cfg.Trace.GRPC.Enabled = endpoint != ""
	}
}

func WithHTTPTracing(endpoint string) Option {
	return func(cfg *Config) {
		cfg.Trace.HTTP.Endpoint = endpoint
		cfg.Trace.HTTP.Enabled = endpoint != ""

	}
}
