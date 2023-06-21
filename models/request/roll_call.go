package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// RollCall describes the `MessageRollCall` message payload.
type RollCall struct {
	From       peer.ID `json:"from,omitempty"`
	Type       string  `json:"type,omitempty"`
	Origin     peer.ID `json:"origin,omitempty"` // Origin is the peer that initiated the roll call.
	FunctionID string  `json:"function_id,omitempty"`
	RequestID  string  `json:"request_id,omitempty"`
}
