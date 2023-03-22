package node

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

// TODO: Potentially move to internal package.
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
		Str("id", peerID.String()).
		Str("addr", maddr.String()).
		Msg("peer connected")

	// Store the peer info.
	err := n.peers.Store(peerID, maddr, addrInfo)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not add peer to peerstore")
	}

	// Update peer list.
	err = n.peers.UpdatePeerList(peerID, maddr, addrInfo)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not update peers in peerstore")
	}
}

func (n *connectionNotifiee) Disconnected(_ network.Network, _ network.Conn) {
	// TBD: Not implemented
}

func (n *connectionNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// TBD: Not implemented
}

func (n *connectionNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// TBD: Not implemented
}
