package node

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO: Check - what do these functions use contexts for?
// TODO: peerID of the sender is a good candidate to move on to the context

type HandlerFunc func(context.Context, peer.ID, []byte) error

func (n *Node) processHealthCheck(ctx context.Context, from peer.ID, payload []byte) error {
	n.log.Debug().Msg("peer health check received")
	return nil
}

func (n *Node) processRollCallResponse(ctx context.Context, from peer.ID, payload []byte) error {
	return errors.New("TBD: Not implemented")
}

func (n *Node) processInstallFunction(ctx context.Context, from peer.ID, payload []byte) error {
	return errors.New("TBD: Not implemented")
}

func (n *Node) processInstallFunctionResponse(ctx context.Context, from peer.ID, payload []byte) error {
	return errors.New("TBD: Not implemented")
}
