package response

import (
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/libp2p/go-libp2p/core/peer"
)

// FormCluster describes the `MessageFormClusteRr` response.
type FormCluster struct {
	Type      string     `json:"type,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
	From      peer.ID    `json:"from,omitempty"`
	Code      codes.Code `json:"code,omitempty"`
}
