package host

import (
	"github.com/blocklessnetworking/b7s/src/models"
)

// defaultConfig used to create Host.
var defaultConfig = Config{
	PrivateKey:          "",
	ConnectionThreshold: 20,
	BootNodes:           nil,
	DialBackPeersLimit:  100,
}

// Config represents the Host configuration.
type Config struct {
	PrivateKey          string
	ConnectionThreshold uint
	BootNodes           []string
	DialBackPeers       []models.Peer
	DialBackPeersLimit  uint
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
func WithBootNodes(nodes []string) func(*Config) {
	return func(cfg *Config) {
		cfg.BootNodes = nodes
	}
}

// WithDialBackPeers specifies dial-back peers that the host initially tries to connect to.
func WithDialBackPeers(peers []models.Peer) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackPeers = peers
	}
}

// WithDialBackPeersLimit specifies the maximum number of dial-back peers to use.
func WithDialBackPeersLimit(n uint) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackPeersLimit = n
	}
}
