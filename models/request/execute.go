package request

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
)

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	From  peer.ID `json:"from,omitempty"`
	Code  string  `json:"code,omitempty"`
	Topic string  `json:"topic,omitempty"`

	execute.Request // execute request is embedded.

	// RequestID may be set initially, if the execution request is relayed via roll-call.
	RequestID string `json:"request_id,omitempty"`

	// Execution request timestamp is a factor for PBFT.
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func (Execute) Type() string { return blockless.MessageExecute }

func (e Execute) MarshalJSON() ([]byte, error) {
	type Alias Execute
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(e),
		Type:  e.Type(),
	}
	return json.Marshal(rec)
}
