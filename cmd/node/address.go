package main

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"
)

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
