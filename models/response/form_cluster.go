package response

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
)

var _ (json.Marshaler) = (*FormCluster)(nil)

// FormCluster describes the `MessageFormClusteRr` response.
type FormCluster struct {
	bls.BaseMessage
	RequestID string         `json:"request_id,omitempty"`
	Code      codes.Code     `json:"code,omitempty"`
	Consensus consensus.Type `json:"consensus,omitempty"`
}

func (f *FormCluster) WithConsensus(c consensus.Type) *FormCluster {
	f.Consensus = c
	return f
}

func (FormCluster) Type() string { return bls.MessageFormClusterResponse }

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
