package node

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_Health(t *testing.T) {

	const (
		testTimeLimit  = 10 * time.Second
		healthInterval = 100 * time.Millisecond
		topic          = DefaultTopic

		expectedPingCount = 3
	)

	var (
		logger          = mocks.NoopLogger
		peerstore       = mocks.BaselinePeerStore(t)
		functionHandler = mocks.BaselineFStore(t)
	)

	// Create a node with a short health interval that will issue quick pings.
	// Then we'll create a host to subscribe to the same topic and verify a few pings before cancelling.

	nhost, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	node, err := New(logger, nhost, peerstore, functionHandler, WithRole(blockless.HeadNode), WithHealthInterval(healthInterval), WithTopic(topic))
	require.NoError(t, err)

	// Create a host that will listen on the the topic to verify health pings
	receiver, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	// Add a deadline for the test so we don't hang.
	ctx, cancel := context.WithTimeout(context.Background(), testTimeLimit)
	defer cancel()

	// Establish a connection between node and receiver.
	hostAddNewPeer(t, nhost, receiver)
	info := hostGetAddrInfo(t, receiver)

	err = node.host.Connect(ctx, *info)
	require.NoError(t, err)

	// Have both client and node subscribe to the same topic.
	_, subscription, err := receiver.Subscribe(ctx, topic)
	require.NoError(t, err)

	_, err = node.subscribe(ctx)
	require.NoError(t, err)

	go node.HealthPing(ctx)

	time.Sleep(subscriptionDiseminationPause)

	// Wait for subscribed messages and verify a few pings came in.
	for i := 0; i < expectedPingCount; i++ {
		msg, err := subscription.Next(ctx)
		require.NoError(t, err)

		require.Equal(t, node.host.ID(), msg.ReceivedFrom)

		var received response.Health
		err = json.Unmarshal(msg.Data, &received)
		require.NoError(t, err)

		require.Equal(t, blockless.MessageHealthCheck, received.Type)
		require.Equal(t, http.StatusOK, received.Code)
	}

	cancel()

	<-ctx.Done()

	// Test should complete but not because of a timeout
	require.NotErrorIsf(t, ctx.Err(), context.DeadlineExceeded, "health test timed out")
}
