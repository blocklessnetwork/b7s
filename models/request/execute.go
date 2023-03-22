package request

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	Type       string              `json:"type,omitempty"`
	From       peer.ID             `json:"from,omitempty"`
	Code       string              `json:"code,omitempty"`
	FunctionID string              `json:"function_id,omitempty"`
	Method     string              `json:"method,omitempty"`
	Parameters []execute.Parameter `json:"parameters,omitempty"`
	Config     execute.Config      `json:"config,omitempty"`

	// RequestID may be set initially, if the execution request is relayed via roll-call.
	RequestID string `json:"request_id,omitempty"`
}
