package pbft

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Option can be used to set PBFT configuration options.
type Option func(*Config)

// PostProcessFunc is invoked by the replica after execution is done.
type PostProcessFunc func(requestID string, origin peer.ID, request execute.Request, result execute.Result)

var DefaultConfig = Config{
	NetworkTimeout: NetworkTimeout,
	RequestTimeout: RequestTimeout,
}

type Config struct {
	PostProcessors []PostProcessFunc // Callback functions to be invoked after execution is done.
	NetworkTimeout time.Duration
	RequestTimeout time.Duration
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
