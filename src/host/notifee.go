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
	// Get the peer ID
	peerID := connection.RemotePeer()
	multiAddr := connection.RemoteMultiaddr()
	peerRecord, err := db.Get(n.Ctx, peerID.Pretty())
	if err != nil {
		log.Info("error getting peer record from database: %v", err)
	}

	peersRecordString, err := db.Get(n.Ctx, "peers")
	if err != nil {
		peersRecordString = []byte("[]")
	}

	var peers []models.Peer
	err = json.Unmarshal(peersRecordString, &peers)
	if err != nil {
		log.Info("error unmarshalling peers record: %v", err)
	}

	peerInfo := models.Peer{
		Type:      "peer",
		Id:        peerID,
		MultiAddr: multiAddr.String(),
		AddrInfo:  network.Peerstore().PeerInfo(peerID),
	}

	j, err := json.Marshal(peerInfo)
	if err != nil {
		log.Info("error marshalling peer info: %v", err)
		return
	}

	peers = append(peers, peerInfo)
	peersRecordString, err = json.Marshal(peers)
	if err != nil {
		log.Info("error marshalling peers record: %v", err)
	}

	log.Info(string(j))
	if peerRecord == nil {
		if err := db.Set(n.Ctx, peerID.Pretty(), string(j)); err != nil {
			log.Info("error setting peer record in database: %v", err)
		}
	}

	log.Info("setting peers in database")
	if err := db.Set(n.Ctx, "peers", string(peersRecordString)); err != nil {
		log.Info("error setting peers record in database: %v", err)
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
