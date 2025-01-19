package bls

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Peer identifies another node in the Bless network.
type Peer struct {
	ID        peer.ID       `json:"id,omitempty"`
	MultiAddr string        `json:"multiaddress,omitempty"`
	AddrInfo  peer.AddrInfo `json:"addrinfo,omitempty"`
}

// PeerIDsToStr will convert a list of peer.IDs to strings.
func PeerIDsToStr(ids []peer.ID) []string {

	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, id.String())
	}

	return out
}
