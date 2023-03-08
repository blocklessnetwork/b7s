package host

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (h *Host) DiscoverPeers(ctx context.Context, topic string) error {

	// Initialize DHT.
	dht, err := h.initDHT(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize DHT: %w", err)
	}

	discovery := routing.NewRoutingDiscovery(dht)
	util.Advertise(ctx, discovery, topic)

	h.log.Debug().Msg("host started peer discovery")

	connected := uint(0)
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

			// Skip peers we're already connected to.
			connections := h.Network().ConnsToPeer(peer.ID)
			if len(connections) > 0 {
				h.log.Debug().
					Str("peer", peer.String()).
					Msg("skipping connected peer")
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
			if connected >= h.cfg.ConnectionThreshold {
				break findPeers
			}
		}

		time.Sleep(h.cfg.DiscoveryInterval)
	}

	h.log.Info().Msg("peer discovery complete")
	return nil
}

func (h *Host) initDHT(ctx context.Context) (*dht.IpfsDHT, error) {

	// Start a DHT for use in peer discovery. Set the DHT to server mode.
	kademliaDHT, err := dht.New(ctx, h.Host, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, fmt.Errorf("could not create DHT: %w", err)
	}

	// Bootstrap the DHT.
	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not bootstrap the DHT: %w", err)
	}

	bootNodes := h.cfg.BootNodes
	peers := h.cfg.DialBackPeers

	// Add the dial-back peers to the list of bootstrap nodes if they do not already exist.
	// TODO: Limit to workers.

	// We may want to limit the number of dial back peers we use.
	added := uint(0)
	addLimit := h.cfg.DialBackPeersLimit

	for _, peer := range peers {
		peer := peer

		// If the limit of dial-back peers is set and we've reached it - stop now.
		if addLimit != 0 && added >= addLimit {
			h.log.Info().Uint("limit", addLimit).Msg("reached limit for dial-back peers")
			break
		}

		addr := peer.String()
		addr = strings.Replace(addr, "127.0.0.1", "0.0.0.0", 1)

		// Check if the peer is already among the boot nodes.
		exists := false
		for _, bootNode := range bootNodes {
			if bootNode.String() == addr {
				exists = true
				break
			}
		}

		// If the peer is not among the boot nodes - add it now.
		if !exists {
			bootNodes = append(bootNodes, peer)
			added++
		}
	}

	// Connect to the bootstrap nodes.
	var wg sync.WaitGroup
	for _, bootNode := range bootNodes {

		addrInfo, err := peer.AddrInfoFromP2pAddr(bootNode)
		if err != nil {
			h.log.Warn().
				Err(err).
				Str("address", bootNode.String()).
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

	return kademliaDHT, nil
}
