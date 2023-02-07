package host

import (
	"github.com/libp2p/go-libp2p/core/network"
)

// Notify sets the notifiee for the underlying libp2p host.
func (h *Host) Notify(notifiee network.Notifiee) {
	h.host.Network().Notify(notifiee)
}
