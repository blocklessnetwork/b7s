package node

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

type connectionNotifiee struct {
	log    zerolog.Logger
	store  blockless.PeerStore
	tracer *tracing.Tracer
}

func newConnectionNotifee(log zerolog.Logger, store blockless.PeerStore) *connectionNotifiee {

	cn := connectionNotifiee{
		log:    log.With().Str("component", "notifiee").Logger(),
		store:  store,
		tracer: tracing.NewTracer("b7s.Notifiee"),
	}

	return &cn
}

func (n *connectionNotifiee) Connected(network network.Network, conn network.Conn) {

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			b7ssemconv.PeerID.String(conn.RemotePeer().String()),
			b7ssemconv.PeerMultiaddr.String(conn.RemoteMultiaddr().String()),
			b7ssemconv.LocalMultiaddr.String(conn.LocalMultiaddr().String()),
		),
	}
	_, span := n.tracer.Start(context.Background(), spanPeerConnected, opts...)
	defer span.End()

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

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			b7ssemconv.PeerID.String(conn.RemotePeer().String()),
			b7ssemconv.PeerMultiaddr.String(conn.RemoteMultiaddr().String()),
			b7ssemconv.LocalMultiaddr.String(conn.LocalMultiaddr().String()),
		),
	}
	_, span := n.tracer.Start(context.Background(), spanPeerDisconnected, opts...)
	defer span.End()

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
