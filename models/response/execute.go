package response

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Execute describes the `MessageExecuteResponse` message payload.
type Execute struct {
	Type      string  `json:"type,omitempty"`
	RequestID string  `json:"request_id,omitempty"`
	From      peer.ID `json:"from,omitempty"`
	Code      string  `json:"code,omitempty"`
	Result    string  `json:"result,omitempty"`
}
