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
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

// TODO: Add an environment variable to skip tests with `publish` due to synchronization issue.

func TestNode_RollCall(t *testing.T) {

	const (
		// TODO: Create a package level const for address and use it in all host creations.
		address = "127.0.0.1"
		port    = 0
	)

	var (
		rollCallReq = request.RollCall{
			Type:       blockless.MessageRollCall,
			FunctionID: "dummy-function-id",
			RequestID:  mocks.GenericUUID.String(),
		}
	)

	payload := serialize(t, rollCallReq)

	t.Run("head node handles roll call", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.HeadNode)
		err := node.processRollCall(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)
	})

	t.Run("worker node handles roll call", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		addr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, addr)

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, blockless.MessageRollCallResponse, received.Type)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, response.CodeAccepted, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("worker node handles failure to check function store", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		addr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, addr)

		// Function store fails to check function presence.
		fstore := mocks.BaselineFunctionHandler(t)
		fstore.GetFunc = func(string, string, bool) (*blockless.FunctionManifest, error) {
			return nil, mocks.GenericError
		}
		node.function = fstore

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, blockless.MessageRollCallResponse, received.Type)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, response.CodeError, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), payload)
		require.Error(t, err)

		wg.Wait()
	})
	t.Run("worker node handles function not installed", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		addr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, addr)

		// Function store works okay but function is not found.
		fstore := mocks.BaselineFunctionHandler(t)
		fstore.GetFunc = func(string, string, bool) (*blockless.FunctionManifest, error) {
			return nil, blockless.ErrNotFound
		}
		node.function = fstore

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			var received response.RollCall
			getStreamPayload(t, stream, &received)

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			require.Equal(t, blockless.MessageRollCallResponse, received.Type)

			require.Equal(t, rollCallReq.FunctionID, received.FunctionID)
			require.Equal(t, rollCallReq.RequestID, received.RequestID)
			require.Equal(t, response.CodeNotFound, received.Code)
		})

		err = node.processRollCall(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("node issues roll call ok", func(t *testing.T) {
		t.Parallel()

		// TODO: Make publishing tests disabled by default and make timeouts longer.

		const (
			topic      = DefaultTopic
			functionID = "super-secret-function-id"
		)

		ctx := context.Background()

		node := createNode(t, blockless.WorkerNode)

		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		addr := getHostAddr(t, receiver)
		info := getAddrInfo(t, addr)
		addPeerToPeerStore(t, node.host, addr)

		err = node.host.Connect(ctx, *info)
		require.NoError(t, err)

		// Have both client and node subscribe to the same topic.
		_, subscription, err := receiver.Subscribe(ctx, topic)
		require.NoError(t, err)

		_, err = node.subscribe(ctx)
		require.NoError(t, err)

		// TODO: Think about how to best handle this.
		time.Sleep(subscriptionDiseminationPause)

		requestID, err := node.issueRollCall(ctx, functionID)
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

		require.Equal(t, blockless.MessageRollCall, received.Type)
		require.Equal(t, functionID, received.FunctionID)
		require.Equal(t, requestID, received.RequestID)
	})
}
