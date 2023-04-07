package response

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the `MessageExecuteResponse` message payload.
type Execute struct {
	Type      string  `json:"type,omitempty"`
	RequestID string  `json:"request_id,omitempty"`
	From      peer.ID `json:"from,omitempty"`
	Code      string  `json:"code,omitempty"`

	// Result is kept for now for backwards compatiblity. It should be
	// equivalent to the `ResultEx.Stdout` field.
	Result   string                `json:"result,omitempty"`
	ResultEx execute.RuntimeOutput `json:"result_ex,omitempty"`
}
