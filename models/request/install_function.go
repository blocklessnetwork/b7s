package request

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

var _ (json.Marshaler) = (*InstallFunction)(nil)

// InstallFunction describes the `MessageInstallFunction` request payload.
type InstallFunction struct {
	From        peer.ID `json:"from,omitempty"`
	ManifestURL string  `json:"manifest_url,omitempty"`
	CID         string  `json:"cid,omitempty"`
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
