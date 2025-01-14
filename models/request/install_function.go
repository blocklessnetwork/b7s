package request

import (
	"encoding/json"
	"errors"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/response"
)

var _ (json.Marshaler) = (*InstallFunction)(nil)

// InstallFunction describes the `MessageInstallFunction` request payload.
type InstallFunction struct {
	bls.BaseMessage
	ManifestURL string `json:"manifest_url,omitempty"`
	CID         string `json:"cid,omitempty"`
}

func (f InstallFunction) Response(c codes.Code) *response.InstallFunction {
	return &response.InstallFunction{
		BaseMessage: bls.BaseMessage{TraceInfo: f.TraceInfo},
		Code:        c,
		Message:     "installed",
		CID:         f.CID,
	}
}

func (InstallFunction) Type() string { return bls.MessageInstallFunction }

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
