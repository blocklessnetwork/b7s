package host

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

// Host represents a new libp2p host.
type Host struct {
	host host.Host
}

// New creates a new Host.
func New(address string, port uint, options ...func(*Config)) (*Host, error) {

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	hostAddress := fmt.Sprintf("/ip4/%v/tcp/%v", address, port)
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(hostAddress),
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

	// Create libp2p host.
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p host: %w", err)
	}

	host := Host{
		host: h,
	}

	return &host, nil
}

// IDs returns the list of p2p IDs of the host.
func (h *Host) IDs() []string {

	// TODO: Perhaps skip local ID..?

	addrs := h.host.Addrs()
	ids := make([]string, 0, len(addrs))

	hostID := h.host.ID()

	for _, addr := range addrs {
		id := fmt.Sprintf("%s/p2p/%s", addr.String(), hostID)
		ids = append(ids, id)
	}

	return ids
}

// readPrivateKey from a file.
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
