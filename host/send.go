package host

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// SendMessage sends a message directly to the specified peer, on the standard blockless protocol.
func (h *Host) SendMessage(ctx context.Context, to peer.ID, payload []byte) error {
	return h.SendMessageOnProtocol(ctx, to, payload, blockless.ProtocolID)
}

// SendMessageOnProtocol sends a message directly to the specified peer, using the specified protocol.
func (h *Host) SendMessageOnProtocol(ctx context.Context, to peer.ID, payload []byte, protocol protocol.ID) error {

	stream, err := h.Host.NewStream(ctx, to, protocol)
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
