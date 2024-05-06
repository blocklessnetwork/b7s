package host

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"
	webtransport "github.com/libp2p/go-libp2p/p2p/transport/webtransport"
	ma "github.com/multiformats/go-multiaddr"
)

// Host represents a new libp2p host.
type Host struct {
	host.Host

	log zerolog.Logger
	cfg Config

	pubsub *pubsub.PubSub
}

// New creates a new Host.
func New(log zerolog.Logger, address string, port uint, options ...func(*Config)) (*Host, error) {
	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	hostAddress := fmt.Sprintf("/ip4/%v/tcp/%v", address, port)
	addresses := []string{hostAddress}

	// define a subset of the default transports, so that we can offer a x509 certificate for the websocket transport
	DefaultTransports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(quic.NewTransport),
		libp2p.Transport(webtransport.New),
	)

	opts := []libp2p.Option{
		DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// Read private key, if provided.
	var key crypto.PrivKey
	var err error

	if cfg.PrivateKey != "" {
		key, err = readPrivateKey(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("could not read private key: %w", err)
		}

		opts = append(opts, libp2p.Identity(key))
	}

	var tlsConfig *tls.Config
	if cfg.Websocket {

		// If the TCP and websocket port are explicitly chosen and set to the same value, one of the two listens will silently fail.
		if port == cfg.WebsocketPort && cfg.WebsocketPort != 0 {
			return nil, fmt.Errorf("TCP and websocket ports cannot be the same (TCP: %v, Websocket: %v)", port, cfg.WebsocketPort)
		}

		// Convert libp2p private key to crypto.PrivateKey
		cryptoPrivKey, err := convertLibp2pPrivKeyToCryptoPrivKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to convert libp2p private key: %v", err)
		}

		// Generate the X.509 certificate
		tlsCert, err := generateX509Certificate(cryptoPrivKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate TLS certificate: %v", err)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			MinVersion:   tls.VersionTLS12,
		}

		wsAddr := fmt.Sprintf("/ip4/%v/tcp/%v/wss", address, cfg.WebsocketPort)
		addresses = append(addresses, wsAddr)
		opts = append(opts, libp2p.Transport(ws.New, ws.WithTLSConfig(tlsConfig)))
	}

	opts = append(opts, libp2p.ListenAddrStrings(addresses...))

	if cfg.DialBackAddress != "" && cfg.DialBackPort != 0 {

		protocol, dialbackAddress, err := determineAddressProtocol(cfg.DialBackAddress)
		if err != nil {
			return nil, fmt.Errorf("could not parse dialback multiaddress (address: %s): %w", cfg.DialBackAddress, err)
		}

		externalAddr := fmt.Sprintf("/%v/%v/tcp/%v", protocol, dialbackAddress, cfg.DialBackPort)
		extAddresses := []string{
			externalAddr,
		}

		if cfg.Websocket && cfg.DialBackWebsocketPort != 0 {

			if cfg.DialBackWebsocketPort == cfg.DialBackPort {
				return nil, fmt.Errorf("TCP and websocket dialback ports cannot be the same (TCP: %v, Websocket: %v)", cfg.DialBackPort, cfg.DialBackWebsocketPort)
			}

			externalWsAddr := fmt.Sprintf("/%v/%v/tcp/%v/ws", protocol, dialbackAddress, cfg.WebsocketPort)
			extAddresses = append(extAddresses, externalWsAddr)
		}

		// Create list of multiaddrs with the external IP and port.
		var externalAddrs []ma.Multiaddr
		for _, addr := range extAddresses {
			maddr, err := ma.NewMultiaddr(addr)
			if err != nil {
				return nil, fmt.Errorf("could not create external multiaddress: %w", err)
			}

			externalAddrs = append(externalAddrs, maddr)
		}

		addrFactory := libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
			// Return only the external multiaddrs.
			return externalAddrs
		})

		opts = append(opts, addrFactory)
	}

	// Create libp2p host.
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p host: %w", err)
	}

	host := Host{
		log: log.With().Str("component", "host").Logger(),
		cfg: cfg,
	}
	host.Host = h

	return &host, nil
}

// PrivateKey returns the private key of the libp2p host.
func (h *Host) PrivateKey() crypto.PrivKey {
	return h.Peerstore().PrivKey(h.ID())
}

// PublicKey returns the public key of the libp2p host.
func (h *Host) PublicKey() crypto.PubKey {
	return h.Peerstore().PubKey(h.ID())
}

// Addresses returns the list of p2p addresses of the host.
func (h *Host) Addresses() []string {

	addrs := h.Addrs()
	out := make([]string, 0, len(addrs))

	hostID := h.ID()

	for _, addr := range addrs {
		addr := fmt.Sprintf("%s/p2p/%s", addr.String(), hostID.String())
		out = append(out, addr)
	}

	return out
}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}

// determineAddressProtocol parses the provided address and tries to determine its type. We typically expect either a IPv4, IPv6 or a hostname.
// At times it's a bit tricky to determine the address type in Go and a lot of parsers end up guessing when dealing with some more exotic variants.
func determineAddressProtocol(address string) (string, string, error) {

	if govalidator.IsIPv4(address) {
		return "ip4", address, nil
	}

	if govalidator.IsIPv6(address) {
		return "ip6", address, nil
	}

	if govalidator.IsDNSName(address) {
		return "dns", address, nil
	}

	return "", "", errors.New("could not parse address")
}
