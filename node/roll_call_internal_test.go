package node

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_RollCall(t *testing.T) {

	t.Run("head node handles roll call", func(t *testing.T) {
		t.Parallel()

		rollCallReq := request.RollCall{
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
		}

		node := createNode(t, blockless.HeadNode)
		err := node.processRollCall(context.Background(), mocks.GenericPeerID, serialize(t, rollCallReq))
		require.NoError(t, err)
	})

	t.Run("worker node handles roll call", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		rollCallReq := request.RollCall{
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
			Origin:     receiver.ID(),
		}

		hostAddNewPeer(t, node.host, receiver)

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, codes.Accepted, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), serialize(t, rollCallReq))
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("worker node handles failure to check function store", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		rollCallReq := request.RollCall{
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
			Origin:     receiver.ID(),
		}

		hostAddNewPeer(t, node.host, receiver)

		// Function store fails to check function presence.
		fstore := mocks.BaselineFStore(t)
		fstore.InstalledFunc = func(string) (bool, error) {
			return false, mocks.GenericError
		}
		node.fstore = fstore

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, codes.Error, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), serialize(t, rollCallReq))
		require.Error(t, err)

		wg.Wait()
	})
	t.Run("worker node installs function on roll call", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		rollCallReq := request.RollCall{
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
			Origin:     receiver.ID(),
		}

		hostAddNewPeer(t, node.host, receiver)

		// Function store has no function but is able to install it.
		fstore := mocks.BaselineFStore(t)
		fstore.InstalledFunc = func(string) (bool, error) {
			return false, nil
		}
		fstore.InstallFunc = func(string, string) error {
			return nil
		}
		node.fstore = fstore

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, codes.Accepted, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), serialize(t, rollCallReq))
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("worker node handles function failure to install function on roll call", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		rollCallReq := request.RollCall{
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
			Origin:     receiver.ID(),
		}

		hostAddNewPeer(t, node.host, receiver)

		// Function store has no function but is not able to install it.
		fstore := mocks.BaselineFStore(t)
		fstore.InstalledFunc = func(string) (bool, error) {
			return false, nil
		}
		fstore.InstallFunc = func(string, string) error {
			return mocks.GenericError
		}
		node.fstore = fstore

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, codes.Error, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), serialize(t, rollCallReq))
		require.Error(t, err)

		wg.Wait()
	})
	t.Run("node issues roll call ok", func(t *testing.T) {
		t.Parallel()

		const (
			topic      = DefaultTopic
			functionID = "super-secret-function-id"
		)

		ctx := context.Background()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		err = receiver.InitPubSub(ctx)
		require.NoError(t, err)

		hostAddNewPeer(t, node.host, receiver)

		info := hostGetAddrInfo(t, receiver)
		err = node.host.Connect(ctx, *info)
		require.NoError(t, err)

		// Have both client and node subscribe to the same topic.
		_, subscription, err := receiver.Subscribe(topic)
		require.NoError(t, err)

		err = node.subscribeToTopics(ctx)
		require.NoError(t, err)

		time.Sleep(subscriptionDiseminationPause)

		requestID, err := newRequestID()
		require.NoError(t, err)

		err = node.publishRollCall(ctx, requestID, functionID, consensus.Type(0), "", nil)
		require.NoError(t, err)

		deadlineCtx, cancel := context.WithTimeout(ctx, publishTimeout)
		defer cancel()

		msg, err := subscription.Next(deadlineCtx)
		require.NoError(t, err)

		from := msg.ReceivedFrom
		require.Equal(t, node.host.ID(), from)
		require.NotNil(t, msg.Topic)
		require.Equal(t, topic, *msg.Topic)

		var received request.RollCall
		err = json.Unmarshal(msg.Data, &received)
		require.NoError(t, err)

		require.Equal(t, functionID, received.FunctionID)
		require.Equal(t, requestID, received.RequestID)
	})
}
