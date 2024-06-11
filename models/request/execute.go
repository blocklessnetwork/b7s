package request

import (
	"encoding/json"
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
)

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	blockless.BaseMessage

	execute.Request // execute request is embedded.

	Topic     string    `json:"topic,omitempty"`
	RequestID string    `json:"request_id,omitempty"` // RequestID may be set initially, if the execution request is relayed via roll-call.
	Timestamp time.Time `json:"timestamp,omitempty"`  // Execution request timestamp is a factor for PBFT.
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
