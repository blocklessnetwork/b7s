package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Execute describes the `MessageExecute` request payload.
type RollCall struct {
	From       peer.ID `json:"from,omitempty"`
	Type       string  `json:"type,omitempty"`
	FunctionID string  `json:"functionId,omitempty"`
	RequestID  string  `json:"request_id,omitempty"`
}
