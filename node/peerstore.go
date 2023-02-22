package node

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type PeerStore interface {
	Store(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
	UpdatePeerList(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error
}
