package request

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/models/bls"
)

var _ (json.Marshaler) = (*DisbandCluster)(nil)

// DisbandCluster describes the `MessageDisbandCluster` request payload.
// It is sent after head node receives the leaders execution response.
type DisbandCluster struct {
	bls.BaseMessage
	RequestID string `json:"request_id,omitempty"`
}

func (DisbandCluster) Type() string { return bls.MessageDisbandCluster }

func (d DisbandCluster) MarshalJSON() ([]byte, error) {
	type Alias DisbandCluster
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(d),
		Type:  d.Type(),
	}
	return json.Marshal(rec)
}
