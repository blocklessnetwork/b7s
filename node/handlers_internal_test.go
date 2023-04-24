package node

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_Handlers(t *testing.T) {

	node := createNode(t, blockless.HeadNode)

	t.Run("health check", func(t *testing.T) {
		t.Parallel()

		msg := response.Health{
			Type: blockless.MessageHealthCheck,
			Code: http.StatusOK,
		}

		payload := serialize(t, msg)
		err := node.processHealthCheck(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)
	})
	t.Run("roll call response", func(t *testing.T) {
		t.Parallel()

		const (
			requestID = "dummy-request-id"
		)

		res := response.RollCall{
			Type:       blockless.MessageRollCallResponse,
			Code:       codes.Accepted,
			Role:       "dummy-role",
			FunctionID: "dummy-function-id",
			RequestID:  requestID,
		}

		// Record response asynchronously.
		var wg sync.WaitGroup
		var recordedResponse response.RollCall
		go func() {
			defer wg.Done()
			recordedResponse = <-node.rollCall.responses(requestID)
		}()

		wg.Add(1)

		payload := serialize(t, res)
		err := node.processRollCallResponse(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)

		wg.Wait()

		expected := res
		expected.From = mocks.GenericPeerID
		require.Equal(t, expected, recordedResponse)
	})
	t.Run("function install response", func(t *testing.T) {
		t.Parallel()

		msg := response.InstallFunction{
			Type:    blockless.MessageInstallFunctionResponse,
			Code:    codes.OK,
			Message: "dummy-message",
		}

		payload := serialize(t, msg)
		err := node.processInstallFunctionResponse(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)
	})
}

func TestNode_InstallFunction(t *testing.T) {

	const (
		manifestURL = "https://example.com/manifest-url"
		cid         = "dummy-cid"
	)

	installReq := request.InstallFunction{
		ManifestURL: manifestURL,
		CID:         cid,
	}

	payload := serialize(t, installReq)

	t.Run("head node handles install", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.HeadNode)

		err := node.processInstallFunction(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)
	})
	t.Run("worker node handles install", func(t *testing.T) {
		t.Parallel()

		const (
			expectedMessage = "installed"
		)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)
		hostAddNewPeer(t, node.host, receiver)

		var wg sync.WaitGroup

		wg.Add(1)
		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received response.InstallFunction
			getStreamPayload(t, stream, &received)

			require.Equal(t, blockless.MessageInstallFunctionResponse, received.Type)
			require.Equal(t, codes.Accepted, received.Code)
			require.Equal(t, expectedMessage, received.Message)
		})

		err = node.processInstallFunction(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("worker node handles function install error", func(t *testing.T) {
		t.Parallel()

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)
		hostAddNewPeer(t, node.host, receiver)

		fstore := mocks.BaselineFunctionHandler(t)
		fstore.InstalledFunc = func(string) (bool, error) {
			return false, nil
		}
		fstore.InstallFunc = func(string, string) error {
			return mocks.GenericError
		}
		node.fstore = fstore

		// NOTE: In reality, this is more "documenting" current behavior.
		// In reality it sounds more correct that we *should* get a response back.
		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			require.Fail(t, "unexpected response")
		})

		err = node.processInstallFunction(context.Background(), receiver.ID(), payload)
		require.Error(t, err)
	})
	t.Run("worker node handles invalid function install requeset", func(t *testing.T) {
		t.Parallel()

		const (
			// JSON without closing brace.
			brokenPayload = `{
				"type": "MsgInstallFunction",
				"manifest_url": "https://example.com/manifest-url",
				"cid": "dummy-cid"`
		)

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)
		hostAddNewPeer(t, node.host, receiver)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			require.Fail(t, "unexpected response")
		})

		err = node.processInstallFunction(context.Background(), receiver.ID(), []byte(brokenPayload))
		require.Error(t, err)
	})
	t.Run("worker node handles failure to send response", func(t *testing.T) {
		t.Parallel()

		// Receiver exists but not added to peer store - the node doesn't know
		// the receivers addresses so `send` will fail.
		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)

		err = node.processInstallFunction(context.Background(), receiver.ID(), payload)
		require.Error(t, err)
	})
}
