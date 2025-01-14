package request

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/response"
)

var _ (json.Marshaler) = (*RollCall)(nil)

// RollCall describes the `MessageRollCall` message payload.
type RollCall struct {
	blockless.BaseMessage
	FunctionID string              `json:"function_id,omitempty"`
	RequestID  string              `json:"request_id,omitempty"`
	Consensus  consensus.Type      `json:"consensus"`
	Attributes *execute.Attributes `json:"attributes,omitempty"`
}

func (r RollCall) Response(c codes.Code) *response.RollCall {
	return &response.RollCall{
		BaseMessage: blockless.BaseMessage{TraceInfo: r.TraceInfo},
		FunctionID:  r.FunctionID,
		RequestID:   r.RequestID,
		Code:        c,
	}
}

func (RollCall) Type() string { return blockless.MessageRollCall }

func (r RollCall) MarshalJSON() ([]byte, error) {
	type Alias RollCall
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(r),
		Type:  r.Type(),
	}
	return json.Marshal(rec)
}
