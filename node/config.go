package node

import (
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

const (
	tracerName = "b7s.Node"

	allowErrorLeakToTelemetry = false // By default we will not send processing errors to telemetry tracers.
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node core.
var DefaultConfig = Config{
	Topics:         []string{blockless.DefaultTopic},
	HealthInterval: blockless.DefaultHealthInterval,
	Concurrency:    blockless.DefaultConcurrency,
}

type Config struct {
	Topics         []string      // Topics to subscribe to.
	HealthInterval time.Duration // How often should we emit the health ping.
	Concurrency    uint          // How many requests should the node process in parallel.
}

// HealthInterval specifies how often we should emit the health signal.
func HealthInterval(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.HealthInterval = d
	}
}

// Concurrency specifies how many requests the node should process in parallel.
func Concurrency(n uint) Option {
	return func(cfg *Config) {
		cfg.Concurrency = n
	}
}

// Topics specifies the p2p topics to which node should subscribe.
func Topics(topics []string) Option {
	return func(cfg *Config) {
		cfg.Topics = topics
	}
}
