package request

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/response"
)

var _ (json.Marshaler) = (*FormCluster)(nil)

// FormCluster describes the `MessageFormCluster` request payload.
// It is sent on clustered execution of a request.
type FormCluster struct {
	bls.BaseMessage
	RequestID      string          `json:"request_id,omitempty"`
	Peers          []peer.ID       `json:"peers,omitempty"`
	Consensus      consensus.Type  `json:"consensus,omitempty"`
	ConnectionInfo []peer.AddrInfo `json:"connection_info,omitempty"`
}

func (f FormCluster) Response(c codes.Code) *response.FormCluster {
	return &response.FormCluster{
		BaseMessage: bls.BaseMessage{TraceInfo: f.TraceInfo},
		RequestID:   f.RequestID,
		Code:        c,
	}

}
func (FormCluster) Type() string { return bls.MessageFormCluster }

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
