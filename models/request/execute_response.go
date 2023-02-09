package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// ExecuteResponse describes the `MessageExecuteResponse` request payload.
type ExecuteResponse struct {
	Type      string  `json:"type,omitempty"`
	RequestID string  `json:"requestId,omitempty"`
	From      peer.ID `json:"from,omitempty"`
	Code      string  `json:"code,omitempty"`
	Result    string  `json:"result,omitempty"`
}
