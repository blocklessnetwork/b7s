package response

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// RollCall describes the `MessageRollCall` response payload.
type RollCall struct {
	Type       string  `json:"type,omitempty"`
	From       peer.ID `json:"from,omitempty"`
	Code       string  `json:"code,omitempty"`
	Role       string  `json:"role,omitempty"`
	FunctionID string  `json:"function_id,omitempty"`
	RequestID  string  `json:"request_id,omitempty"`
}
