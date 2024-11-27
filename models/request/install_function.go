package request

import (
	"encoding/json"
	"errors"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

var _ (json.Marshaler) = (*InstallFunction)(nil)

// InstallFunction describes the `MessageInstallFunction` request payload.
type InstallFunction struct {
	blockless.BaseMessage
	ManifestURL string `json:"manifest_url,omitempty"`
	CID         string `json:"cid,omitempty"`
}

func (f InstallFunction) Response(c codes.Code) *response.InstallFunction {
	return &response.InstallFunction{
		BaseMessage: blockless.BaseMessage{TraceInfo: f.TraceInfo},
		Code:        c,
		Message:     "installed",
		CID:         f.CID,
	}
}

func (InstallFunction) Type() string { return blockless.MessageInstallFunction }

func (f InstallFunction) MarshalJSON() ([]byte, error) {
	type Alias InstallFunction
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(f),
		Type:  f.Type(),
	}
	return json.Marshal(rec)
}

func (f InstallFunction) Valid() error {

	if f.CID == "" {
		return errors.New("function CID is required")
	}

	return nil
}
