package response

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
)

// InstallFunction describes the response to the `MessageInstallFunction` message.
type InstallFunction struct {
	Type    string     `json:"type,omitempty"`
	From    peer.ID    `json:"from,omitempty"`
	Code    codes.Code `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
}
