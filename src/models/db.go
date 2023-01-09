package models

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	Type      string        `json:"type,omitempty"`
	Id        peer.ID       `json:"id,omitempty"`
	MultiAddr string        `json:"multiaddress,omitempty"`
	AddrInfo  peer.AddrInfo `json:"addrinfo,omitempty"`
}
