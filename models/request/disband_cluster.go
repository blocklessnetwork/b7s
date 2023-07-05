package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// DisbandCluster describes the `MessageDisbandCluster` request payload.
// It is sent after head node receives the leaders execution response.
type DisbandCluster struct {
	Type      string  `json:"type,omitempty"`
	From      peer.ID `json:"from,omitempty"`
	RequestID string  `json:"request_id,omitempty"`
}
