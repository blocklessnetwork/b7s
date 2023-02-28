package node

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_Messaging(t *testing.T) {

	const (
		clientAddress = "127.0.0.1"
		clientPort    = 0

		// TODO: Use a different topic.
		topic = DefaultTopic
	)

	var (
		rec = dummyRecord{
			ID:          mocks.GenericUUID.String(),
			Value:       19846,
			Description: "dummy-description",
		}
	)

	client, err := host.New(mocks.NoopLogger, clientAddress, clientPort)
	require.NoError(t, err)

	addr := getHostAddr(t, client)

	node := createNode(t, blockless.HeadNode)
	addPeerToPeerStore(t, node.host, addr)

	t.Run("sending single message", func(t *testing.T) {
		t.Parallel()

		var wg sync.WaitGroup
		wg.Add(1)

		client.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received dummyRecord
			getStreamPayload(t, stream, &received)

			require.Equal(t, rec, received)
		})

		err = node.send(context.Background(), client.ID(), rec)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("publishing to a topic", func(t *testing.T) {
		t.Parallel()

		const (
			// How long can the client wait for a published message before giving up.
			clientTimeout = 3 * time.Second

			// It seems like a delay is needed so that the hosts exchange information about the fact
			// that they are subscribed to the same topic. If that does not happen, node might publish
			// a message too soon and the client might miss it.
			// It will then wait for a published message in vain.
			// This is the pause we make after subscribing to the topic and before publishing a message.
			subscriptionrDiseminationPause = 250 * time.Millisecond
		)

		ctx := context.Background()

		// Establish a connection between peers.
		clientInfo := getAddrInfo(t, addr)
		err = node.host.Connect(ctx, *clientInfo)
		require.NoError(t, err)

		// Have both client and node subscribe to the same topic.
		_, subscription, err := client.Subscribe(ctx, topic)
		require.NoError(t, err)

		_, err = node.subscribe(ctx)
		require.NoError(t, err)

		// TODO: Think about how to best handle this.
		time.Sleep(subscriptionrDiseminationPause)

		err = node.publish(ctx, rec)
		require.NoError(t, err)

		deadlineCtx, cancel := context.WithTimeout(ctx, clientTimeout)
		defer cancel()
		msg, err := subscription.Next(deadlineCtx)
		require.NoError(t, err)

		from := msg.ReceivedFrom
		require.Equal(t, node.host.ID(), from)
		require.NotNil(t, msg.Topic)
		require.Equal(t, topic, *msg.Topic)

		var received dummyRecord
		err = json.Unmarshal(msg.Data, &received)
		require.NoError(t, err)

		require.Equal(t, rec, received)
	})
}

type dummyRecord struct {
	ID          string `json:"id"`
	Value       uint64 `json:"value"`
	Description string `json:"description"`
}
