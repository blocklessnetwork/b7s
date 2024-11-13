package node

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_SendMessage(t *testing.T) {

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	node := createNode(t, blockless.HeadNode)
	hostAddNewPeer(t, node.host, client)

	rec := newDummyRecord()

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

	err = node.send(context.Background(), client.ID(), &rec)
	require.NoError(t, err)

	wg.Wait()
}

func TestNode_Publish(t *testing.T) {

	var (
		rec   = newDummyRecord()
		ctx   = context.Background()
		topic = DefaultTopic
	)

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	err = client.InitPubSub(ctx)
	require.NoError(t, err)

	node := createNode(t, blockless.HeadNode)
	hostAddNewPeer(t, node.host, client)

	err = node.subscribeToTopics(ctx)
	require.NoError(t, err)

	// Establish a connection between peers.
	clientInfo := hostGetAddrInfo(t, client)
	err = node.host.Connect(ctx, *clientInfo)
	require.NoError(t, err)

	// Have both client and node subscribe to the same topic.
	_, subscription, err := client.Subscribe(topic)
	require.NoError(t, err)

	time.Sleep(subscriptionDiseminationPause)

	err = node.publish(ctx, &rec)
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
}

func TestNode_SendMessageToMany(t *testing.T) {

	client1, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	client2, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	node := createNode(t, blockless.HeadNode)
	hostAddNewPeer(t, node.host, client1)
	hostAddNewPeer(t, node.host, client2)

	client1.SetStreamHandler(blockless.ProtocolID, func(network.Stream) {})
	client2.SetStreamHandler(blockless.ProtocolID, func(network.Stream) {})

	// NOTE: These subtests are sequential.
	t.Run("nominal case - sending to two online peers is ok", func(t *testing.T) {
		err = node.sendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), true)
		require.NoError(t, err)
	})
	t.Run("peer is down with requireAll is an error", func(t *testing.T) {
		client1.Close()
		err = node.sendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), true)
		require.Error(t, err)
	})
	t.Run("peer is down with partial delivery is ok", func(t *testing.T) {
		client1.Close()
		err = node.sendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), false)
		require.NoError(t, err)
	})
	t.Run("all sends failing produces an error", func(t *testing.T) {
		client1.Close()
		client2.Close()
		err = node.sendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), false)
		require.Error(t, err)
	})
}

type dummyRecord struct {
	ID          string `json:"id"`
	Value       uint64 `json:"value"`
	Description string `json:"description"`
}

func (dummyRecord) Type() string {
	return "MessageDummyRecord"
}

func newDummyRecord() dummyRecord {
	return dummyRecord{
		ID:          mocks.GenericUUID.String(),
		Value:       rand.Uint64(),
		Description: "dummy-description",
	}
}
