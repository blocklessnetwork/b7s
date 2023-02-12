package response

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// InstallFunction describes the response to the `MessageInstallFunction` message.
type InstallFunction struct {
	Type    string  `json:"type,omitempty"`
	From    peer.ID `json:"from,omitempty"`
	Code    string  `json:"code,omitempty"`
	Message string  `json:"message,omitempty"`
}
