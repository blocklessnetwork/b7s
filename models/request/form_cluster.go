package request

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
)

var _ (json.Marshaler) = (*FormCluster)(nil)

// FormCluster describes the `MessageFormCluster` request payload.
// It is sent on clustered execution of a request.
type FormCluster struct {
	RequestID string         `json:"request_id,omitempty"`
	Peers     []peer.ID      `json:"peers,omitempty"`
	Consensus consensus.Type `json:"consensus,omitempty"`
}

func (FormCluster) Type() string { return blockless.MessageFormCluster }

func (f FormCluster) MarshalJSON() ([]byte, error) {
	type Alias FormCluster
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(f),
		Type:  f.Type(),
	}
	return json.Marshal(rec)
}
