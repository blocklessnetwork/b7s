package host

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

// Host represents a new libp2p host.
type Host struct {
	log       zerolog.Logger
	host.Host // TODO: Once the use cases cristalize - reconsider embedding vs private field

	cfg Config
}

// New creates a new Host.
func New(log zerolog.Logger, address string, port uint, options ...func(*Config)) (*Host, error) {

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	hostAddress := fmt.Sprintf("/ip4/%v/tcp/%v", address, port)
	addresses := []string{
		hostAddress,
	}

	if cfg.Websocket {
		wsAddr := hostAddress + "/ws"
		addresses = append(addresses, wsAddr)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(addresses...),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// Read private key, if provided.
	if cfg.PrivateKey != "" {
		key, err := readPrivateKey(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("could not read private key: %w", err)
		}

		opts = append(opts, libp2p.Identity(key))
	}

	if cfg.DialBackAddress != "" && cfg.DialBackPort != 0 {

		externalAddr := fmt.Sprintf("/ip4/%s/tcp/%d", cfg.DialBackAddress, cfg.DialBackPort)
		addrs := []string{
			externalAddr,
		}

		if cfg.Websocket {
			externalWsAddr := externalAddr + "/ws"
			addresses = append(addrs, externalWsAddr)
		}

		// Create list of multiaddrs with the external IP and port.
		var externalAddrs []ma.Multiaddr
		for _, addr := range addresses {
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
