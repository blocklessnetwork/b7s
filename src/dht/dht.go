package dht

import (
	"context"
	"sync"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

func InitDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)

	// all nodes should respond to queries
	dht.Mode(dht.ModeServer)

	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	bootNodes := []multiaddr.Multiaddr{}

	cfg := ctx.Value("config").(models.Config)
	for _, bootNode := range cfg.Node.BootNodes {
		bootNodes = append(bootNodes, multiaddr.StringCast(bootNode))
	}

	for _, peerAddr := range bootNodes {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				// todo figure out what we want to do with no good addresses
				// no reason to panic here get's noisy with discovery
				if err.Error() != "no good addresses" {
					log.WithFields(log.Fields{
						"localMultiAddr": h.Addrs(),
						"peerID":         h.ID(),
						"err":            err,
					}).Warn("bootstrap warn")
				}
			}
		}()
	}

	wg.Wait()

	return kademliaDHT
}

func DiscoverPeers(ctx context.Context, h host.Host, topicName string) {
	kademliaDHT := InitDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, topicName)
	log.Info("starting peer discovery")
	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {

		peerChan, err := routingDiscovery.FindPeers(ctx, topicName)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peer)
			if err != nil {
				// this can be quite noisy with discovery
				// fmt.Println("Failed connecting to ", peer.ID.Pretty(), ", error:", err)
			} else {
				pebble := db.Get(peer.ID.Pretty())
				db.Set(pebble, "peerID", peer.ID.Pretty())
				log.WithFields(log.Fields{
					"peerID": peer.ID.Pretty(),
				}).Info("connected to peer")
				anyConnected = true
			}
		}
	}
	log.Info("Peer discovery complete")
}
