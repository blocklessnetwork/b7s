package response

import (
	"encoding/json"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
)

var _ (json.Marshaler) = (*FormCluster)(nil)

// FormCluster describes the `MessageFormClusteRr` response.
type FormCluster struct {
	blockless.BaseMessage
	RequestID string         `json:"request_id,omitempty"`
	Code      codes.Code     `json:"code,omitempty"`
	Consensus consensus.Type `json:"consensus,omitempty"`
}

func (FormCluster) Type() string { return blockless.MessageFormClusterResponse }

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
