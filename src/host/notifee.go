package host

import (
	"context"
	"log"

	db "github.com/blocklessnetworking/b7s/src/db"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

type ConnectedNotifee struct {
	Ctx context.Context
}

// Implement the Connected/Disconnected methods of the Notifee interface
func (n *ConnectedNotifee) Connected(network network.Network, connection network.Conn) {

	// get the peer id
	peerID := connection.RemotePeer()
	// get the multiaddress
	multiAddr := connection.RemoteMultiaddr()
	peerRecord, err := db.Get(n.Ctx, peerID.Pretty())

	if err != nil {
		log.Println(err)
	}

	if peerRecord == nil {
		// peer is not in the database, add it
		db.Set(n.Ctx, peerID.Pretty(), multiAddr.String())
	}

	log.Println("Connected to: ", multiAddr, peerID)
}

func (n *ConnectedNotifee) Disconnected(network network.Network, connection network.Conn) {
	// A peer has been disconnected
	// Do something with the disconnected peer
}

func (n *ConnectedNotifee) Listen(network.Network, ma.Multiaddr) {
	// A new stream has been opened
	// Do something with the stream
}

func (n *ConnectedNotifee) ListenClose(network.Network, ma.Multiaddr) {
	// A stream has been closed
	// Do something with the closed stream
}
