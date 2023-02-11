package host

import (
	"context"
	"fmt"
	"strings"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetworking/b7s/src/models"
)

// TODO: bootNodes and peers.
func (h *Host) DiscoverPeers(ctx context.Context, topic string, bootNodes []string, peers []models.Peer) error {

	// Initialize DHT.
	dht, err := h.initDHT(ctx, bootNodes, peers)
	if err != nil {
		return fmt.Errorf("could not initalize DHT: %w", err)
	}

	discovery := routing.NewRoutingDiscovery(dht)
	util.Advertise(ctx, discovery, topic)

	h.log.Debug().Msg("host started peer discovery")

	connected := 0
findPeers:
	for {
		peers, err := discovery.FindPeers(ctx, topic)
		if err != nil {
			return fmt.Errorf("could not find peers: %w", err)
		}

		for peer := range peers {
			// Skip self.
			if peer.ID == h.ID() {
				continue
			}

			err = h.Connect(ctx, peer)
			if err != nil {
				h.log.Debug().
					Err(err).
					Str("peer", peer.String()).
					Msg("could not connect to peer")
				continue
			}

			h.log.Info().Str("peer", peer.String()).Msg("connected to peer")

			connected++

			// Stop when we have reached connection threshold.
			if connected >= connectionThreshold {
				break findPeers
			}
		}
	}

	h.log.Info().Msg("peer discovery complete")
	return nil
}

func (h *Host) initDHT(ctx context.Context, bootNodes []string, peers []models.Peer) (*dht.IpfsDHT, error) {

	// Start a DHT for use in peer discovery.
	kademlieDHT, err := dht.New(ctx, h.Host)
	if err != nil {
		return nil, fmt.Errorf("could not create DHT: %w", err)
	}

	// Set the DHT to server mode.
	dht.Mode(dht.ModeServer)

	// Bootstrap the DHT.
	err = kademlieDHT.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not bootstrap the DHT: %w", err)
	}

	// Add the dial-back peers to the list of bootrstrap nodes if they do not already exist.
	// TODO: Limit the number of dial-back peers.
	// TODO: Limit to workers.
	for _, peer := range peers {

		addr := fmt.Sprintf("%s/p2p/%s", peer.MultiAddr, peer.Id.String())
		addr = strings.Replace(addr, "127.0.0.1", "0.0.0.0", 1)

		// Check if the peer is already among the boot nodes.
		exists := false
		for _, bootNode := range bootNodes {
			if bootNode == addr {
				exists = true
				break
			}
		}

		// If it's not - add it now.
		if !exists {
			bootNodes = append(bootNodes, addr)
		}
	}

	// Connect to the bootstrap nodes.
	var wg sync.WaitGroup
	for _, bootNode := range bootNodes {

		maddr, err := multiaddr.NewMultiaddr(bootNode)
		if err != nil {
			h.log.Warn().
				Err(err).
				Str("address", bootNode).
				Msg("could not parse multiaddress for boot node - skipping")
			continue
		}

		addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			h.log.Warn().
				Err(err).
				Str("address", bootNode).
				Msg("could not get addrinfo for boot node - skipping")
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			peerAddr := addrInfo

			err := h.Host.Connect(ctx, *peerAddr)
			if err != nil {
				if err.Error() != errNoGoodAddresses {
					h.log.Error().
						Err(err).
						Str("addrinfo", peerAddr.String()).
						Msg("could not connect to bootstrap node")
				}
			}
		}()
	}

	// Wait until we know the outcome of all connection attempts.
	wg.Wait()

	return kademlieDHT, nil
}
