package node

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

type connectionNotifiee struct {
	log   zerolog.Logger
	peers PeerStore
}

func newConnectionNotifee(log zerolog.Logger, peerStore PeerStore) *connectionNotifiee {

	cn := connectionNotifiee{
		log:   log.With().Str("component", "notifiee").Logger(),
		peers: peerStore,
	}

	return &cn
}

func (n *connectionNotifiee) Connected(network network.Network, conn network.Conn) {

	// Get peer information.
	peerID := conn.RemotePeer()
	maddr := conn.RemoteMultiaddr()
	addrInfo := network.Peerstore().PeerInfo(peerID)

	n.log.Debug().
		Str("peer", peerID.String()).
		Str("addr", maddr.String()).
		Msg("peer connected")

	// Store the peer info.
	err := n.peers.Store(peerID, maddr, addrInfo)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not add peer to peerstore")
	}
}

func (n *connectionNotifiee) Disconnected(_ network.Network, conn network.Conn) {

	// TODO: Check - do we want to remove peer after he's been disconnected.

	peerID := conn.RemotePeer()
	n.log.Debug().
		Str("peer", peerID.String()).
		Msg("peer disconnected")
}

func (n *connectionNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func (n *connectionNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}
