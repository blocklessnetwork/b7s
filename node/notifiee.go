package node

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type connectionNotifiee struct {
	log   zerolog.Logger
	store blockless.PeerStore
}

func newConnectionNotifee(log zerolog.Logger, store blockless.PeerStore) *connectionNotifiee {

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

	// We could save only the mutliaddress from which we receive this connection. However, we could theoretically have multiple connections
	// and there's no reason to limit ourselves to a single address.

	peer := blockless.Peer{
		ID:        peerID,
		MultiAddr: maddr.String(),
		// AddrInfo struct basically repeats the above info (multiaddress).
		AddrInfo: peer.AddrInfo{
			ID:    peerID,
			Addrs: make([]multiaddr.Multiaddr, 0),
		},
	}

	for _, conn := range network.ConnsToPeer(conn.RemotePeer()) {
		peer.AddrInfo.Addrs = append(peer.AddrInfo.Addrs, conn.RemoteMultiaddr())
	}

	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_address", maddr.String()).
		Str("local_address", laddr.String()).
		Any("addr_info", peer.AddrInfo).
		Msg("peer connected")

	// Store the peer info.
	err := n.store.SavePeer(peer)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not add peer to peerstore")
	}
}

func (n *connectionNotifiee) Disconnected(_ network.Network, conn network.Conn) {

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
