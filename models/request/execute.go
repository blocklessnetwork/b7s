package request

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
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

func (e Execute) Response(c codes.Code) *response.Execute {
	return &response.Execute{
		BaseMessage: blockless.BaseMessage{TraceInfo: e.TraceInfo},
		RequestID:   e.RequestID,
		Code:        c,
	}
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

func (e Execute) Valid() error {

	c, err := consensus.Parse(e.Config.ConsensusAlgorithm)
	if err != nil {
		return fmt.Errorf("could not parse consensus algorithm: %w", err)
	}

	if c == consensus.PBFT &&
		e.Config.NodeCount > 0 &&
		e.Config.NodeCount < pbft.MinimumReplicaCount {
		return fmt.Errorf("minimum %v nodes needed for PBFT consensus", pbft.MinimumReplicaCount)
	}

	return nil
}
