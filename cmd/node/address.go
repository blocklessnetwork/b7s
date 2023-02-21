package main

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// getPeerAddresses returns the list of the multiaddreses for the peer list.
func getPeerAddresses(peers []blockless.Peer) ([]multiaddr.Multiaddr, error) {

	var addrs []multiaddr.Multiaddr

	for _, peer := range peers {

		addr, err := multiaddr.NewMultiaddr(peer.MultiAddr)
		if err != nil {
			return nil, fmt.Errorf("could not parse multiaddress (addr: %s): %w", peer.MultiAddr, err)
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// parse list of strings with multiaddresses
func getBootNodeAddresses(addrs []string) ([]multiaddr.Multiaddr, error) {

	var out []multiaddr.Multiaddr
	for _, addr := range addrs {

		addr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("could not parse multiaddress (addr: %s): %w", addr, err)
		}

		out = append(out, addr)
	}

	return out, nil
}
