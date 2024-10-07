package response

import (
	"encoding/json"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the response to the `MessageExecute` message.
type Execute struct {
	blockless.BaseMessage
	RequestID string            `json:"request_id,omitempty"`
	Code      codes.Code        `json:"code,omitempty"`
	Results   execute.ResultMap `json:"results,omitempty"`
	Cluster   execute.Cluster   `json:"cluster,omitempty"`

	// Used to communicate the reason for failure to the user.
	Message string `json:"message,omitempty"`
}

func (e *Execute) WithResults(r execute.ResultMap) *Execute {
	e.Results = r
	return e
}

func (e *Execute) WithCluster(c execute.Cluster) *Execute {
	e.Cluster = c
	return e
}

func (Execute) Type() string { return blockless.MessageExecuteResponse }

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
