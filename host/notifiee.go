package host

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blessnetwork/b7s/telemetry/tracing"
)

type Notifiee struct {
	log    zerolog.Logger
	store  bls.PeerStore
	tracer *tracing.Tracer
}

func NewNotifee(log zerolog.Logger, store bls.PeerStore) *Notifiee {

	cn := Notifiee{
		log:    log.With().Str("component", "notifiee").Logger(),
		store:  store,
		tracer: tracing.NewTracer("b7s.Notifiee"),
	}

	return &cn
}

func (n *Notifiee) Connected(network network.Network, conn network.Conn) {

	ctx, span := n.tracer.Start(context.Background(), spanPeerConnected, connectionTraceOpts(conn)...)
	defer span.End()

	// Get peer information.
	peerID := conn.RemotePeer()
	maddr := conn.RemoteMultiaddr()
	laddr := conn.LocalMultiaddr()

	peer := bls.Peer{
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
	err := n.store.SavePeer(ctx, peer)
	if err != nil {
		n.log.Warn().Err(err).Str("id", peerID.String()).Msg("could not add peer to peerstore")
	}
}

func (n *Notifiee) Disconnected(_ network.Network, conn network.Conn) {

	_, span := n.tracer.Start(context.Background(), spanPeerDisconnected, connectionTraceOpts(conn)...)
	defer span.End()

	maddr := conn.RemoteMultiaddr()
	laddr := conn.LocalMultiaddr()

	peerID := conn.RemotePeer()
	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_address", maddr.String()).
		Str("local_address", laddr.String()).
		Msg("peer disconnected")
}

func (n *Notifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func (n *Notifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func connectionTraceOpts(conn network.Conn) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithAttributes(
			b7ssemconv.PeerID.String(conn.RemotePeer().String()),
			b7ssemconv.PeerMultiaddr.String(conn.RemoteMultiaddr().String()),
			b7ssemconv.LocalMultiaddr.String(conn.LocalMultiaddr().String()),
		),
	}
}
