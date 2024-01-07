package node

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_Messaging(t *testing.T) {

	const (
		topic = DefaultTopic
	)

	var (
		rec = dummyRecord{
			ID:          mocks.GenericUUID.String(),
			Value:       19846,
			Description: "dummy-description",
		}
	)

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	node := createNode(t, blockless.HeadNode)
	hostAddNewPeer(t, node.host, client)

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

		err := node.send(context.Background(), client.ID(), rec)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("publishing to a topic", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		err = client.InitPubSub(ctx)
		require.NoError(t, err)

		// Establish a connection between peers.
		clientInfo := hostGetAddrInfo(t, client)
		err = node.host.Connect(ctx, *clientInfo)
		require.NoError(t, err)

		// Have both client and node subscribe to the same topic.
		_, subscription, err := client.Subscribe(topic)
		require.NoError(t, err)

		err = node.subscribeToTopics(ctx)
		require.NoError(t, err)

		time.Sleep(subscriptionDiseminationPause)

		err = node.publish(ctx, rec)
		require.NoError(t, err)

		deadlineCtx, cancel := context.WithTimeout(ctx, publishTimeout)
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
