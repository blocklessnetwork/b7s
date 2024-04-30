package node

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type connectionNotifiee struct {
	log   zerolog.Logger
	store Store
}

func newConnectionNotifee(log zerolog.Logger, store Store) *connectionNotifiee {

	cn := connectionNotifiee{
		log:   log.With().Str("component", "notifiee").Logger(),
		store: store,
	}

	return &cn
}

func (n *connectionNotifiee) Connected(network network.Network, conn network.Conn) {

	// Get peer information.
	peerID := conn.RemotePeer()
	maddr := conn.RemoteMultiaddr()
	laddr := conn.LocalMultiaddr()
	addrInfo := network.Peerstore().PeerInfo(peerID)

	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_address", maddr.String()).
		Str("local_address", laddr.String()).
		Interface("addr_info", addrInfo).
		Msg("peer connected")

	peer := blockless.Peer{
		ID:        peerID,
		MultiAddr: maddr.String(),
		AddrInfo:  addrInfo,
	}

	// Store the peer info.
	err := n.store.SavePeer(peer)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not add peer to peerstore")
	}
}

func (n *connectionNotifiee) Disconnected(_ network.Network, conn network.Conn) {

	// TODO: Check - do we want to remove peer after he's been disconnected.
	maddr := conn.RemoteMultiaddr()
	laddr := conn.LocalMultiaddr()

	peerID := conn.RemotePeer()
	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_address", maddr.String()).
		Str("local_address", laddr.String()).
		Msg("peer disconnected")
}

func (n *connectionNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func (n *connectionNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}
