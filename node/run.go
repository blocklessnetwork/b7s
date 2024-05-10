package node

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// Run will start the main loop for the node.
func (n *Node) Run(ctx context.Context) error {

	err := n.subscribeToTopics(ctx)
	if err != nil {
		return fmt.Errorf("could not subscribe to topics: %w", err)
	}

	// Sync functions now in case they were removed from the storage.
	err = n.fstore.Sync(false)
	if err != nil {
		return fmt.Errorf("could not sync functions: %w", err)
	}

	// Set the handler for direct messages.
	n.listenDirectMessages(ctx)

	// Discover peers.
	// NOTE: Potentially signal any error here so that we abort the node
	// run loop if anything failed.
	for _, topic := range n.cfg.Topics {
		go func(topic string) {

			// TODO: Check DHT initialization, now that we're working with multiple topics, may not need to repeat ALL work per topic.
			err = n.host.DiscoverPeers(ctx, topic)
			if err != nil {
				n.log.Error().Err(err).Msg("could not discover peers")
			}

		}(topic)
	}

	// Start the health signal emitter in a separate goroutine.
	go n.HealthPing(ctx)

	// Start the function sync in the background to periodically check functions.
	go n.runSyncLoop(ctx)

	n.log.Info().Uint("concurrency", n.cfg.Concurrency).Msg("starting node main loop")

	var workers sync.WaitGroup

	// Process topic messages - spin up a goroutine for each topic that will feed the main processing loop below.
	// No need for locking since we're still single threaded here and these (subscribed) topics will not be touched by other code.
	for name, topic := range n.subgroups.topics {

		workers.Add(1)

		go func(name string, subscription *pubsub.Subscription) {
			defer workers.Done()

			// Message processing loops.
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

				n.log.Trace().Str("topic", name).Str("peer", msg.ReceivedFrom.String()).Str("id", msg.ID).Msg("received message")

				// Try to get a slot for processing the request.
				n.sema <- struct{}{}
				n.wg.Add(1)

				go func(msg *pubsub.Message) {
					// Free up slot after we're done.
					defer n.wg.Done()
					defer func() { <-n.sema }()

					err = n.processMessage(ctx, msg.ReceivedFrom, msg.GetData(), subscriptionPipeline)
					if err != nil {
						n.log.Error().Err(err).Str("id", msg.ID).Str("peer", msg.ReceivedFrom.String()).Msg("could not process message")
					}
				}(msg)
			}
		}(name, topic.subscription)
	}

	workers.Wait()

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

		n.log.Trace().Str("peer", from.String()).Msg("received direct message")

		err = n.processMessage(ctx, from, msg, directMessagePipeline)
		if err != nil {
			n.log.Error().Err(err).Str("peer", from.String()).Msg("could not process direct message")
		}
	})
}
