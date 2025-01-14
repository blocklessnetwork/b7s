package node_test

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/host"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/node"
	"github.com/blessnetwork/b7s/testing/helpers"
	"github.com/blessnetwork/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"

	// It seems like a delay is needed so that the hosts exchange information about the fact
	// that they are subscribed to the same topic. If that does not happen, node might publish
	// a message too soon and the client might miss it. It will then wait for a published message in vain.
	// This is the pause we make after subscribing to the topic and before publishing a message.
	// In reality as little as 250ms is enough, but lets allow a longer time for when
	// tests are executed in parallel or on weaker machines.
	subscriptionDiseminationPause = 2 * time.Second

	// How long can the client wait for a published message before giving up.
	publishTimeout = 10 * time.Second
)

func TestNode_Publish(t *testing.T) {

	var (
		rec   = newDummyRecord()
		ctx   = context.Background()
		topic = bls.DefaultTopic
	)

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	err = client.InitPubSub(ctx)
	require.NoError(t, err)

	core := createNodeCore(t)

	err = core.Host().InitPubSub(ctx)
	require.NoError(t, err)

	err = core.Subscribe(ctx, bls.DefaultTopic)
	require.NoError(t, err)

	// Establish a connection between peers.
	clientInfo := helpers.HostGetAddrInfo(t, client)
	err = core.Host().Connect(ctx, *clientInfo)
	require.NoError(t, err)

	// Have both client and node subscribe to the same topic.
	_, subscription, err := client.Subscribe(topic)
	require.NoError(t, err)

	time.Sleep(subscriptionDiseminationPause)

	err = core.Publish(ctx, &rec)
	require.NoError(t, err)

	deadlineCtx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()
	msg, err := subscription.Next(deadlineCtx)
	require.NoError(t, err)

	from := msg.ReceivedFrom
	require.Equal(t, core.Host().ID(), from)
	require.NotNil(t, msg.Topic)
	require.Equal(t, topic, *msg.Topic)

	var received dummyRecord
	err = json.Unmarshal(msg.Data, &received)
	require.NoError(t, err)
	require.Equal(t, rec, received)
}

func TestNode_SendMessageToMany(t *testing.T) {

	var (
		log = mocks.NoopLogger

		client1 = helpers.NewLoopbackHost(t, log)
		client2 = helpers.NewLoopbackHost(t, log)

		core = node.NewCore(log, helpers.NewLoopbackHost(t, log))
	)

	client1.SetStreamHandler(bls.ProtocolID, func(stream network.Stream) {})
	client2.SetStreamHandler(bls.ProtocolID, func(stream network.Stream) {})

	// NOTE: These subtests are sequential.

	// At this point we don't know how to dial the clients so sends will fail.
	t.Run("all sends failing produces an error", func(t *testing.T) {
		err := core.SendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), false)
		require.Error(t, err)
	})

	// Add client1 to peerstore - now one send can succeed.
	helpers.HostAddNewPeer(t, core.Host(), client1)

	t.Run("peer is down with requireAll is an error", func(t *testing.T) {
		err := core.SendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), true)
		require.Error(t, err)
	})
	t.Run("peer is down with partial delivery is ok", func(t *testing.T) {
		err := core.SendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), false)
		require.NoError(t, err)
	})

	// Add client2 to the peerstore - now both sends can succeed.
	helpers.HostAddNewPeer(t, core.Host(), client2)

	t.Run("nominal case - sending to two online peers is ok", func(t *testing.T) {
		err := core.SendToMany(context.Background(), []peer.ID{client1.ID(), client2.ID()}, newDummyRecord(), true)
		require.NoError(t, err)
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

func createNodeCore(t *testing.T) node.Core {

	logger := mocks.NoopLogger
	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	return node.NewCore(logger, host)
}
