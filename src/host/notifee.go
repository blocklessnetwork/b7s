package host

import (
	"context"
	"encoding/json"

	db "github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// ConnectedNotifee is a struct that implements the Notifee interface
type ConnectedNotifee struct {
	Ctx context.Context
}

func (n *ConnectedNotifee) Connected(network network.Network, connection network.Conn) {

	// get the peer id and multiaddr.
	peerID := connection.RemotePeer()

	// this address doesn't give us the in bound port for a dialback
	multiAddr := connection.RemoteMultiaddr()

	// Get the peer record from the database.
	peerRecord, err := db.Get(n.Ctx, peerID.Pretty())
	if err != nil {
		log.WithError(err).Info("Error getting peer record from database")
	}

	// Get the list of peers from the database.
	peersRecordString, err := db.Get(n.Ctx, "peers")
	if err != nil {
		peersRecordString = []byte("[]")
	}

	var peers []models.Peer
	if err = json.Unmarshal(peersRecordString, &peers); err != nil {
		log.WithError(err).Info("Error unmarshalling peers record")
	}

	// Create a new peer info struct.
	peerInfo := models.Peer{
		Type:      "peer",
		Id:        peerID,
		MultiAddr: multiAddr.String(),
		AddrInfo:  network.Peerstore().PeerInfo(peerID),
	}

	// Check if the peer is already in the list.
	peerExists := false
	for _, peer := range peers {
		if peer.Id == peerInfo.Id {
			peerExists = true
			break
		}
	}

	// If the peer is not in the list, add it.
	if !peerExists {
		peers = append(peers, peerInfo)
	}

	// Marshal the peer info struct to JSON.
	peerJSON, err := json.Marshal(peerInfo)
	if err != nil {
		log.WithError(err).Info("Error marshalling peer info")
		return
	}

	// Marshal the list of peers to JSON.
	peersJSON, err := json.Marshal(peers)
	if err != nil {
		log.WithError(err).Info("Error marshalling peers record")
	}

	//log the peerInfo
	log.WithFields(log.Fields{
		"peerInfo": peerInfo,
	}).Info("Peer Info Stored")

	// If the peer record does not exist in the database, set it.
	if peerRecord == nil {
		if err := db.Set(n.Ctx, peerID.Pretty(), string(peerJSON)); err != nil {
			log.WithError(err).Info("Error setting peer record in database")
		}
	}

	// Set the list of peers in the database.
	if err := db.Set(n.Ctx, "peers", string(peersJSON)); err != nil {
		log.WithError(err).Info("Error setting peers record in database")
	}
}

// Disconnected is called when a connection is closed
func (n *ConnectedNotifee) Disconnected(network network.Network, connection network.Conn) {
	// A peer has been disconnected
	// Do something with the disconnected peer
}

// Listen is called when a new stream is opened
func (n *ConnectedNotifee) Listen(network.Network, ma.Multiaddr) {
	// A new stream has been opened
	// Do something with the stream
}

// ListenClose is called when a stream is closed
func (n *ConnectedNotifee) ListenClose(network.Network, ma.Multiaddr) {
	// A stream has been closed
	// Do something with the closed stream
}
