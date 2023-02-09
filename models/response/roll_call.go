package response

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// RollCall describes the `RollCall` response payload.
type RollCall struct {
	Type       string  `json:"type,omitempty"`
	From       peer.ID `json:"from,omitempty"`
	Code       string  `json:"code,omitempty"`
	Role       string  `json:"role,omitempty"`
	FunctionID string  `json:"functionId,omitempty"`
	RequestID  string  `json:"request_id,omitempty"`
}
