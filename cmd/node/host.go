package main

import (
	"fmt"

	"github.com/blocklessnetwork/b7s/config"
	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/rs/zerolog"
)

func createHost(log zerolog.Logger, cfg config.Config, dialbackPeers ...blockless.Peer) (*host.Host, error) {

	// Get the list of boot nodes addresses.
	bootNodes, err := getBootNodeAddresses(cfg.BootNodes)
	if err != nil {
		return nil, fmt.Errorf("could not get boot node addresses: %w", err)
	}

	opts := []func(*host.Config){
		host.WithPrivateKey(cfg.Connectivity.PrivateKey),
		host.WithBootNodes(bootNodes),
		host.WithDialBackAddress(cfg.Connectivity.DialbackAddress),
		host.WithDialBackPort(cfg.Connectivity.DialbackPort),
		host.WithDialBackWebsocketPort(cfg.Connectivity.WebsocketDialbackPort),
		host.WithWebsocket(cfg.Connectivity.Websocket),
		host.WithWebsocketPort(cfg.Connectivity.WebsocketPort),
		host.WithDialBackPeers(dialbackPeers),
	}

	// Create libp2p host.
	host, err := host.New(log, cfg.Connectivity.Address, cfg.Connectivity.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create host (key: '%s'): %w", cfg.Connectivity.PrivateKey, err)
	}
	defer host.Close()

	return host, nil
}
