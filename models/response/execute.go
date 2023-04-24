package response

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the `MessageExecuteResponse` message payload.
type Execute struct {
	Type      string                    `json:"type,omitempty"`
	RequestID string                    `json:"request_id,omitempty"`
	From      peer.ID                   `json:"from,omitempty"`
	Code      codes.Code                `json:"code,omitempty"`
	Results   map[string]execute.Result `json:"results,omitempty"`

	// Used to communicate the reason for failure to the user.
	Message string `json:"message,omitempty"`
}
