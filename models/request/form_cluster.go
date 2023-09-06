package request

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
)

// FormCluster describes the `MessageFormCluster` request payload.
// It is sent on clustered execution of a request.
type FormCluster struct {
	Type      string         `json:"type,omitempty"`
	From      peer.ID        `json:"from,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Peers     []peer.ID      `json:"peers,omitempty"`
	Consensus consensus.Type `json:"consensus,omitempty"`
}
