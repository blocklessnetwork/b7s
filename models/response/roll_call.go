package response

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
)

// RollCall describes the `MessageRollCall` response payload.
type RollCall struct {
	Type       string     `json:"type,omitempty"`
	From       peer.ID    `json:"from,omitempty"`
	Code       codes.Code `json:"code,omitempty"`
	Role       string     `json:"role,omitempty"`
	FunctionID string     `json:"function_id,omitempty"`
	RequestID  string     `json:"request_id,omitempty"`
}
