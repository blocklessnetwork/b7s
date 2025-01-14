package response

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/models/blockless"
)

var _ (json.Marshaler) = (*Health)(nil)

// Health describes the message sent as a health ping.
type Health struct {
	blockless.BaseMessage
	Code int `json:"code,omitempty"`
}

func (Health) Type() string { return blockless.MessageHealthCheck }

func (h Health) MarshalJSON() ([]byte, error) {
	type Alias Health
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(h),
		Type:  h.Type(),
	}
	return json.Marshal(rec)
}
