package request

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/consensus/pbft"
	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/response"
	"github.com/hashicorp/go-multierror"
)

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	blockless.BaseMessage

	execute.Request // execute request is embedded.

	Topic string `json:"topic,omitempty"`
}

func (e Execute) Response(c codes.Code, id string) *response.Execute {
	return &response.Execute{
		BaseMessage: blockless.BaseMessage{TraceInfo: e.TraceInfo},
		RequestID:   id,
		Code:        c,
	}
}

func (e Execute) RollCall(id string, c consensus.Type) *RollCall {
	return &RollCall{
		BaseMessage: blockless.BaseMessage{TraceInfo: e.TraceInfo},
		RequestID:   id,
		FunctionID:  e.FunctionID,
		Consensus:   c,
		Attributes:  e.Config.Attributes,
	}
}

func (e Execute) WorkOrder(id string) *WorkOrder {
	return &WorkOrder{
		BaseMessage: blockless.BaseMessage{TraceInfo: e.TraceInfo},
		RequestID:   id,
		Request:     e.Request,
		Timestamp:   time.Now().UTC(),
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

	var multierr *multierror.Error
	err := e.Request.Valid()
	if err != nil {
		multierr = multierror.Append(multierr, err)
	}

	c, err := consensus.Parse(e.Config.ConsensusAlgorithm)
	if err != nil {
		multierr = multierror.Append(multierr, fmt.Errorf("could not parse consensus algorithm: %w", err))
	}

	if c == consensus.PBFT &&
		e.Config.NodeCount > 0 &&
		e.Config.NodeCount < pbft.MinimumReplicaCount {

		multierr = multierror.Append(multierr, fmt.Errorf("minimum %v nodes needed for PBFT consensus", pbft.MinimumReplicaCount))
	}

	return multierr.ErrorOrNil()
}
