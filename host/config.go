package host

import (
	"time"

	"github.com/multiformats/go-multiaddr"
)

// defaultConfig used to create Host.
var defaultConfig = Config{
	PrivateKey:          "",
	ConnectionThreshold: 20,
	DialBackPeersLimit:  100,
	DiscoveryInterval:   10 * time.Second,
}

// Config represents the Host configuration.
type Config struct {
	PrivateKey          string
	ConnectionThreshold uint
	BootNodes           []multiaddr.Multiaddr
	DialBackPeers       []multiaddr.Multiaddr
	DialBackPeersLimit  uint
	DiscoveryInterval   time.Duration
	DialBackAddress     string
	DialBackPort        uint
}

// WithPrivateKey specifies the private key for the Host.
func WithPrivateKey(filepath string) func(*Config) {
	return func(cfg *Config) {
		cfg.PrivateKey = filepath
	}
}

// WithConnectionThreshold specifies how many connections should the host wait for on peer discovery.
func WithConnectionThreshold(n uint) func(*Config) {
	return func(cfg *Config) {
		cfg.ConnectionThreshold = n
	}
}

// WithBootNodes specifies boot nodes that the host initially tries to connect to.
func WithBootNodes(nodes []multiaddr.Multiaddr) func(*Config) {
	return func(cfg *Config) {
		cfg.BootNodes = nodes
	}
}

// WithDialBackPeers specifies dial-back peers that the host initially tries to connect to.
func WithDialBackPeers(peers []multiaddr.Multiaddr) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackPeers = peers
	}
}

func WithDialBackAddress(a string) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackAddress = a
	}
}

func WithDialBackPort(n uint) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackPort = n
	}
}

// WithDialBackPeersLimit specifies the maximum number of dial-back peers to use.
func WithDialBackPeersLimit(n uint) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackPeersLimit = n
	}
}

// WithDiscoveryInterval specifies how often we should try to discover new peers during the discovery phase.
func WithDiscoveryInterval(d time.Duration) func(*Config) {
	return func(cfg *Config) {
		cfg.DiscoveryInterval = d
	}
}
