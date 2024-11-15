package response

import (
	"encoding/json"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

type WorkOrder struct {
	blockless.BaseMessage

	RequestID string             `json:"request_id,omitempty"`
	Code      codes.Code         `json:"code,omitempty"`
	Result    execute.NodeResult `json:"result,omitempty"`

	ErrorMessage string `json:"error_message,omitempty"`
}

func (w *WorkOrder) WithMetadata(m any) *WorkOrder {
	w.Result.Metadata = m
	return w
}

func (e *WorkOrder) WithErrorMessage(err error) *WorkOrder {
	e.ErrorMessage = err.Error()
	return e
}

func (WorkOrder) Type() string { return blockless.MessageWorkOrderResponse }

func (e WorkOrder) MarshalJSON() ([]byte, error) {
	type Alias WorkOrder
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(e),
		Type:  e.Type(),
	}
	return json.Marshal(rec)
}
