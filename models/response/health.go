package response

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Health describes the message sent as a health ping.
type Health struct {
	Type string  `json:"type,omitempty"`
	From peer.ID `json:"from,omitempty"`
	Code string  `json:"code,omitempty"`
}
