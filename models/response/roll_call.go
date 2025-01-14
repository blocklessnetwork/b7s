package response

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
)

var _ (json.Marshaler) = (*RollCall)(nil)

// RollCall describes the `MessageRollCall` response payload.
type RollCall struct {
	blockless.BaseMessage
	Code       codes.Code `json:"code,omitempty"`
	FunctionID string     `json:"function_id,omitempty"`
	RequestID  string     `json:"request_id,omitempty"`
}

func (RollCall) Type() string { return blockless.MessageRollCallResponse }

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
