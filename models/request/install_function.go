package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// InstallFunction describes the `MessageInstallFunction` request payload.
type InstallFunction struct {
	Type        string  `json:"type,omitempty"`
	From        peer.ID `json:"from,omitempty"`
	ManifestURL string  `json:"manifest_url,omitempty"`
	CID         string  `json:"cid,omitempty"`
}
