package host

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// SendMessage sends a message directly to the specified peer.
func (h *Host) SendMessage(ctx context.Context, to peer.ID, payload []byte) error {

	stream, err := h.Host.NewStream(ctx, to, blockless.ProtocolID)
	if err != nil {
		return fmt.Errorf("could not create stream: %w", err)
	}
	defer stream.Close()

	_, err = stream.Write(payload)
	if err != nil {
		stream.Reset()
		return fmt.Errorf("could not write payload: %w", err)
	}

	return nil
}
