package pbft

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/metadata"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

// Option can be used to set PBFT configuration options.
type Option func(*Config)

// PostProcessFunc is invoked by the replica after execution is done.
type PostProcessFunc func(requestID string, origin peer.ID, request execute.Request, result execute.Result)

var DefaultConfig = Config{
	NetworkTimeout:   NetworkTimeout,
	RequestTimeout:   RequestTimeout,
	MetadataProvider: metadata.NewNoopProvider(),
}

type Config struct {
	PostProcessors   []PostProcessFunc // Callback functions to be invoked after execution is done.
	NetworkTimeout   time.Duration
	RequestTimeout   time.Duration
	MetadataProvider metadata.Provider
	TraceInfo        tracing.TraceInfo
}

// WithNetworkTimeout sets how much time we allow for message sending.
func WithNetworkTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.NetworkTimeout = d
	}
}

// WithRequestTimeout sets the inactivity period before we trigger a view change.
func WithRequestTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.RequestTimeout = d
	}
}

// WithPostProcessors sets the callbacks that will be invoked after execution.
func WithPostProcessors(callbacks ...PostProcessFunc) Option {
	return func(cfg *Config) {
		var fns []PostProcessFunc
		fns = append(fns, callbacks...)
		cfg.PostProcessors = fns
	}
}

// WithMetadataProvider sets the metadata provider for the node.
func WithMetadataProvider(p metadata.Provider) Option {
	return func(cfg *Config) {
		cfg.MetadataProvider = p
	}
}

// WithTraceInfo passes along telemetry trace information.
func WithTraceInfo(t tracing.TraceInfo) Option {
	return func(cfg *Config) {
		cfg.TraceInfo = t
	}
}
