package aggregate

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/execute"
)

type Results []Result

// Result represents the execution result along with its aggregation stats.
type Result struct {
	Result execute.RuntimeOutput `json:"result,omitempty"`
	// Peers that got this result.
	Peers []peer.ID `json:"peers,omitempty"`
	// Peers metadata
	Metadata NodeMetadata `json:"metadata,omitempty"`
	// How frequent was this result, in percentages.
	Frequency float64 `json:"frequency,omitempty"`
}

type NodeMetadata map[peer.ID]any

func (m NodeMetadata) MarshalJSON() ([]byte, error) {

	em := make(map[string]any, len(m))
	for p, v := range m {
		em[p.String()] = v
	}

	return json.Marshal(em)
}
