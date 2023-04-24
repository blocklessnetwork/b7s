package blockless

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Peer identifies another node in the Blockless network.
type Peer struct {
	ID        peer.ID       `json:"id,omitempty"`
	MultiAddr string        `json:"multiaddress,omitempty"`
	AddrInfo  peer.AddrInfo `json:"addrinfo,omitempty"`
}
