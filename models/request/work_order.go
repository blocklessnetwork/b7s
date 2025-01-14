package request

import (
	"encoding/json"
	"time"

	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/response"
)

type WorkOrder struct {
	blockless.BaseMessage

	execute.Request // execute request is embedded

	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"` // Execution request timestamp is a factor for PBFT.
}

func (w WorkOrder) Response(c codes.Code, res execute.Result) *response.WorkOrder {

	return &response.WorkOrder{
		BaseMessage: blockless.BaseMessage{TraceInfo: w.TraceInfo},
		Code:        c,
		RequestID:   w.RequestID,
		Result: execute.NodeResult{
			Result: res,
		},
	}
}

func (WorkOrder) Type() string { return blockless.MessageWorkOrder }

func (w WorkOrder) MarshalJSON() ([]byte, error) {
	type Alias WorkOrder
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(w),
		Type:  w.Type(),
	}
	return json.Marshal(rec)
}
