package request

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

var _ (json.Marshaler) = (*FormCluster)(nil)

// FormCluster describes the `MessageFormCluster` request payload.
// It is sent on clustered execution of a request.
type FormCluster struct {
	blockless.BaseMessage
	RequestID string         `json:"request_id,omitempty"`
	Peers     []peer.ID      `json:"peers,omitempty"`
	Consensus consensus.Type `json:"consensus,omitempty"`
}

func (f FormCluster) Response(c codes.Code) *response.FormCluster {
	return &response.FormCluster{
		BaseMessage: blockless.BaseMessage{TraceInfo: f.TraceInfo},
		RequestID:   f.RequestID,
		Code:        c,
	}

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
