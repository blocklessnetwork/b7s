package node

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Run will start the main loop for the node.
func (n *Node) Run(ctx context.Context) error {

	// Subscribe to the specified topic.
	subscription, err := n.subscribe(ctx)
	if err != nil {
		return fmt.Errorf("could not subscribe to topic: %w", err)
	}

	// Set the handler for direct messages.
	n.listenDirectMessages(ctx)

	// Discover peers.
	// NOTE: Potentially signal any error here so that we abort the node
	// run loop if anything failed.
	go func() {
		err = n.host.DiscoverPeers(ctx, n.topicName)
		if err != nil {
			n.log.Error().
				Err(err).
				Msg("could not discover peers")
		}
	}()

	// Start the health signal emitter in a separate goroutine.
	go n.HealthPing(ctx)

	n.log.Info().Msg("starting node main loop")

	// Message processing loop.
	for {

		// Retrieve next message.
		msg, err := subscription.Next(ctx)
		if err != nil {
			// NOTE: Cancelling the context will lead us here.
			n.log.Error().Err(err).Msg("could not receive message")
			break
		}

		// Skip messages we published.
		if msg.ReceivedFrom == n.host.ID() {
			continue
		}

		n.log.Debug().
			Str("id", msg.ID).
			Str("peer_id", msg.ReceivedFrom.String()).
			Msg("received message")

		err = n.processMessage(ctx, msg.ReceivedFrom, msg.Data)
		if err != nil {
			n.log.Error().
				Err(err).
				Str("id", msg.ID).
				Str("peer_id", msg.ReceivedFrom.String()).
				Msg("could not process message")
			continue
		}
	}

	return nil
}

// listenDirectMessages will process messages sent directly to the peer (as opposed to published messages).
func (n *Node) listenDirectMessages(ctx context.Context) {

	n.host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()

		buf := bufio.NewReader(stream)
		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			n.log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		err = n.processMessage(ctx, from, msg)
		if err != nil {
			n.log.Error().
				Err(err).
				Str("peer_id", from.String()).
				Msg("could not process direct message")
		}
	})
}
