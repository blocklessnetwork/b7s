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

	n.log.Info().Strs("topics", n.cfg.Topics).Msg("topics node will subscribe to")

	err := n.subscribeToTopics(ctx)
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
	// TODO: SUS1 - DHT stuff inside gets multiplied.
	for _, topic := range n.cfg.Topics {
		go func(topic string) {
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

	msgs := make(chan *pubsub.Message, n.cfg.Concurrency)
	var topicWorkers sync.WaitGroup

	// Process topic messages - spin up a goroutine for each topic that will feed the main processing loop below.
	for name, topic := range n.topics {

		topicWorkers.Add(1)

		go func(name string, subscription *pubsub.Subscription) {
			defer topicWorkers.Done()

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

				msgs <- msg
			}
		}(name, topic.subscription)
	}

	// Read and process messages.
	go func() {
		for msg := range msgs {

			n.log.Debug().Str("peer", msg.ReceivedFrom.String()).Str("id", msg.ID).Msg("processing message")

			n.wg.Add(1)
			go func(msg *pubsub.Message) {
				defer n.wg.Done()

				err = n.processMessage(ctx, msg.ReceivedFrom, msg.Data)
				if err != nil {
					n.log.Error().Err(err).Str("id", msg.ID).Str("peer", msg.ReceivedFrom.String()).Msg("could not process message")
				}
			}(msg)
		}
	}()

	// Waiting for topic workers to stop (context canceled).
	topicWorkers.Wait()
	// Signal that no new messages will be incoming.
	close(msgs)

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
