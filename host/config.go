package host

import (
	"time"

	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// defaultConfig used to create Host.
var defaultConfig = Config{
	PrivateKey:                         "",
	ConnectionThreshold:                20,
	DialBackPeersLimit:                 100,
	DiscoveryInterval:                  10 * time.Second,
	Websocket:                          false,
	BootNodesReachabilityCheckInterval: 1 * time.Minute,
	MustReachBootNodes:                 defaultMustReachBootNodes,
	EnableP2PRelay:                     false,
}

// Config represents the Host configuration.
type Config struct {
	PrivateKey string

	ConnectionThreshold uint
	BootNodes           []multiaddr.Multiaddr
	DialBackPeers       []blockless.Peer
	DialBackPeersLimit  uint
	DiscoveryInterval   time.Duration

	Websocket     bool
	WebsocketPort uint

	DialBackAddress       string
	DialBackPort          uint
	DialBackWebsocketPort uint

	BootNodesReachabilityCheckInterval time.Duration
	MustReachBootNodes                 bool
	DisableResourceLimits              bool
	EnableP2PRelay                     bool
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
func WithDialBackPeers(peers []blockless.Peer) func(*Config) {
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

func WithDialBackWebsocketPort(n uint) func(*Config) {
	return func(cfg *Config) {
		cfg.DialBackWebsocketPort = n
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

// WithWebsocket specifies whether libp2p host should use websocket protocol.
func WithWebsocket(b bool) func(*Config) {
	return func(cfg *Config) {
		cfg.Websocket = b
	}
}

// WithWebsocketPort specifies on which port the host should listen for websocket connections.
func WithWebsocketPort(port uint) func(*Config) {
	return func(cfg *Config) {
		cfg.WebsocketPort = port
	}
}

// WithMustReachBootNodes specifies if we should treat failure to reach boot nodes as a halting error.
func WithMustReachBootNodes(b bool) func(*Config) {
	return func(cfg *Config) {
		cfg.MustReachBootNodes = b
	}
}

// WithBootNodesReachabilityInterval specifies how often should we recheck and reconnect to boot nodes.
func WithBootNodesReachabilityInterval(d time.Duration) func(cfg *Config) {
	return func(cfg *Config) {
		cfg.BootNodesReachabilityCheckInterval = d
	}
}

// WithDisabledResourceLimits allows removing any resource limits set by libp2p.
// WARNING: experimental
func WithDisabledResourceLimits(b bool) func(cfg *Config) {
	return func(cfg *Config) {
		cfg.DisableResourceLimits = b
	}
}

// EnableP2PRelay allows user to control whether the b7s can act as a p2p relayer for worker nodes
func WithEnableP2PRelay(b bool) func(cfg *Config) {
	return func(cfg *Config) {
		cfg.EnableP2PRelay = b
	}
}
