package host

import (
	"context"
	"fmt"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetwork/b7s/models/blockless"
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
		h.log.Trace().Msg("starting peer discovery")

		// Using a list instead of a channel. If this starts getting too large switch back.
		peers, err := util.FindPeers(ctx, discovery, topic)
		if err != nil {
			return fmt.Errorf("could not find peers: %w", err)
		}

		h.log.Trace().Int("count", len(peers)).Str("topic", topic).Msg("discovered peers")

		for _, peer := range peers {
			// Skip self.
			if peer.ID == h.ID() {
				continue
			}

			// Skip peers we're already connected to.
			connections := h.Network().ConnsToPeer(peer.ID)
			if len(connections) > 0 {
				continue
			}

			err = h.Connect(ctx, peer)
			if err != nil {
				h.log.Debug().Err(err).Str("peer", peer.ID.String()).Msg("could not connect to discovered peer")
				continue
			}

			h.log.Info().Str("peer", peer.ID.String()).Msg("connected to peer")

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

	// Nodes we will try to connect to on start.
	var bootNodes []blockless.Peer

	// Add explicitly specified nodes first.
	for _, addr := range h.cfg.BootNodes {
		addr := addr

		addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			h.log.Warn().Err(err).Str("address", addr.String()).Msg("could not get addrinfo for boot node - skipping")
			continue
		}

		node := blockless.Peer{
			ID:       addrInfo.ID,
			AddrInfo: *addrInfo,
		}

		bootNodes = append(bootNodes, node)
	}

	// Add the dial-back peers to the list of bootstrap nodes if they do not already exist.

	// We may want to limit the number of dial back peers we use.
	added := uint(0)
	addLimit := h.cfg.DialBackPeersLimit

	var dialbackPeers []blockless.Peer
	for _, peer := range h.cfg.DialBackPeers {
		peer := peer

		// If the limit of dial-back peers is set and we've reached it - stop now.
		if addLimit != 0 && added >= addLimit {
			h.log.Info().Uint("limit", addLimit).Msg("reached limit for dial-back peers")
			break
		}

		// If we don't have any addresses, add the multiaddress we (hopefully) do have - last one we received a connection from.
		if len(peer.AddrInfo.Addrs) == 0 {

			ma, err := multiaddr.NewMultiaddr(peer.MultiAddr)
			if err != nil {
				h.log.Warn().Str("peer", peer.ID.String()).Str("addr", peer.MultiAddr).Msg("invalid multiaddress for dial-back peer, skipping")
				continue
			}

			h.log.Debug().Str("peer", peer.ID.String()).Msg("using last known multiaddress for dial-back peer")

			peer.AddrInfo.Addrs = []multiaddr.Multiaddr{ma}
		}

		h.log.Debug().Str("peer", peer.ID.String()).Interface("addr_info", peer.AddrInfo).Msg("adding dial-back peer")

		dialbackPeers = append(dialbackPeers, peer)
		added++
	}

	bootNodes = append(bootNodes, dialbackPeers...)

	// Connect to the bootstrap nodes.
	var wg sync.WaitGroup
	for _, bootNode := range bootNodes {
		bootNode := bootNode

		// Skip peers we're already connected to (perhaps a dial-back peer was also a boot node).
		connections := h.Network().ConnsToPeer(bootNode.ID)
		if len(connections) > 0 {
			continue
		}

		wg.Add(1)
		go func(peer blockless.Peer) {
			defer wg.Done()

			peerAddr := peer.AddrInfo

			err := h.Host.Connect(ctx, peerAddr)
			if err != nil {
				if err.Error() != errNoGoodAddresses {
					h.log.Error().Err(err).Str("peer", peer.ID.String()).Interface("addr_info", peerAddr).Msg("could not connect to bootstrap node")
				}

				return
			}

			h.log.Debug().Str("peer", peer.ID.String()).Any("addr_info", peerAddr).Msg("connected to known peer")
		}(bootNode)
	}

	// Wait until we know the outcome of all connection attempts.
	wg.Wait()

	return kademliaDHT, nil
}
