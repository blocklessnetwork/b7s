package response

import (
	"encoding/json"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
)

var _ (json.Marshaler) = (*InstallFunction)(nil)

// InstallFunction describes the response to the `MessageInstallFunction` message.
type InstallFunction struct {
	bls.BaseMessage
	Code    codes.Code `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	CID     string     `json:"cid,omitempty"`
}

func (InstallFunction) Type() string { return bls.MessageInstallFunctionResponse }

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
