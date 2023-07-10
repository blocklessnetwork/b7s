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

	// Sync functions now in case they were removed from the storage.
	n.syncFunctions()

	// Set the handler for direct messages.
	n.listenDirectMessages(ctx)

	// Discover peers.
	// NOTE: Potentially signal any error here so that we abort the node
	// run loop if anything failed.
	go func() {
		err = n.host.DiscoverPeers(ctx, n.cfg.Topic)
		if err != nil {
			n.log.Error().Err(err).Msg("could not discover peers")
		}
	}()

	// Start the health signal emitter in a separate goroutine.
	go n.HealthPing(ctx)

	// Start the function sync in the background to periodically check functions.
	go n.runSyncLoop(ctx)

	n.log.Info().Uint("concurrency", n.cfg.Concurrency).Msg("starting node main loop")

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

		n.log.Trace().Str("id", msg.ID).Str("peer", msg.ReceivedFrom.String()).Msg("received message")

		// Try to get a slot for processing the request.
		n.sema <- struct{}{}
		n.wg.Add(1)

		go func() {
			// Free up slot after we're done.
			defer n.wg.Done()
			defer func() { <-n.sema }()

			err = n.processMessage(ctx, msg.ReceivedFrom, msg.Data)
			if err != nil {
				n.log.Error().Err(err).Str("id", msg.ID).Str("peer_id", msg.ReceivedFrom.String()).Msg("could not process message")
			}
		}()
	}

	n.log.Debug().Msg("waiting for messages being processed")

	n.wg.Wait()

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

		n.log.Debug().Str("peer", from.String()).Msg("received direct message")

		err = n.processMessage(ctx, from, msg)
		if err != nil {
			n.log.Error().Err(err).Str("peer", from.String()).Msg("could not process direct message")
		}
	})
}
