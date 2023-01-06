package dht

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	db "github.com/blocklessnetworking/b7s/src/db"
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
	// Start a DHT, for use in peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		log.Fatal(err)
	}

	// Set the DHT to server mode.
	dht.Mode(dht.ModeServer)

	// Bootstrap the DHT.
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatal(err)
	}

	// Get the list of bootstrap nodes from the configuration.
	cfg := ctx.Value("config").(models.Config)
	bootNodes := cfg.Node.BootNodes

	// Get the list of dial-back peers from the database.
	var dialBackPeers []models.Peer
	peersRecordString, err := db.Get(ctx, "peers")
	if err != nil {
		peersRecordString = []byte("[]")
	}
	if err = json.Unmarshal(peersRecordString, &dialBackPeers); err != nil {
		log.WithError(err).Info("Error unmarshalling peers record")
	}

	//log the length of dialBackPeers
	log.WithField("dialBackPeers", len(dialBackPeers)).Info("dialBackPeers")

	// Convert the dial-back peers to multiaddrs and add them to the list of bootstrap nodes if they do not already exist.
	// likely good to limit the number of dial-back peers to a small number.
	// and we need to limit to workers
	for _, peer := range dialBackPeers {
		peerMultiAddr := fmt.Sprintf("%s/p2p/%s", peer.MultiAddr, peer.Id.Pretty())
		peerMultiAddr = strings.Replace(peerMultiAddr, "127.0.0.1", "0.0.0.0", 1)
		//log peer add
		log.WithField("peerMultiAddr", peerMultiAddr).Info("peerMultiAddr")
		peerExists := false
		for _, bootNode := range bootNodes {
			if bootNode == peerMultiAddr {
				peerExists = true
				break
			}
		}
		if !peerExists {
			bootNodes = append(bootNodes, peerMultiAddr)
		}
	}

	// Connect to the bootstrap nodes.
	var wg sync.WaitGroup
	for _, bootNode := range bootNodes {
		peerAddr, err := peer.AddrInfoFromP2pAddr(multiaddr.StringCast(bootNode))
		log.Info("booting from: ", peerAddr)
		if err != nil {
			log.WithFields(log.Fields{
				"bootNode": bootNode,
				"error":    err,
			}).Warn("Invalid bootstrap node address")
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerAddr); err != nil {
				if err.Error() != "no good addresses" {
					log.WithFields(log.Fields{
						"localMultiAddr": h.Addrs(),
						"peerID":         h.ID(),
						"err":            err,
					}).Warn("Error connecting to bootstrap node")
				}
			}
		}()
	}
	wg.Wait()
	return kademliaDHT
}

func DiscoverPeers(ctx context.Context, h host.Host) {
	topicName := ctx.Value("topicName").(string)
	kademliaDHT := InitDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, topicName)
	log.Info("starting peer discovery")
	// Look for others who have announced and attempt to connect to them
	numConnected := 0
	for numConnected < 20 {
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
				log.WithFields(log.Fields{
					"peerID": peer.ID.Pretty(),
				}).Info("connected to a peer")
				numConnected++
				if numConnected >= 20 {
					break
				}
			}
		}
	}
	log.Info("Peer discovery complete")
}
