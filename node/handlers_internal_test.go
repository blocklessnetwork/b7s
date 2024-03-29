package node

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_Handlers(t *testing.T) {

	node := createNode(t, blockless.HeadNode)

	t.Run("health check", func(t *testing.T) {
		t.Parallel()

		msg := response.Health{
			Code: http.StatusOK,
		}

		err := node.processHealthCheck(context.Background(), mocks.GenericPeerID, msg)
		require.NoError(t, err)
	})
	t.Run("roll call response", func(t *testing.T) {
		t.Parallel()

		const (
			requestID = "dummy-request-id"
		)

		node.rollCall.create(requestID)

		res := response.RollCall{
			Code:       codes.Accepted,
			FunctionID: "dummy-function-id",
			RequestID:  requestID,
		}

		// Record response asynchronously.
		var wg sync.WaitGroup
		var recordedResponse rollCallResponse
		go func() {
			defer wg.Done()
			recordedResponse = <-node.rollCall.responses(requestID)
		}()

		wg.Add(1)

		err := node.processRollCallResponse(context.Background(), mocks.GenericPeerID, res)
		require.NoError(t, err)

		wg.Wait()

		expected := res
		require.Equal(t, expected, recordedResponse.RollCall)
	})
	t.Run("skipping inadequate roll call responses", func(t *testing.T) {
		t.Parallel()

		const (
			requestID = "dummy-request-id-2"
		)

		node.rollCall.create(requestID)

		// We only want responses with the code `Accepted`.
		res := response.RollCall{
			Code:       codes.NotFound,
			RequestID:  requestID,
			FunctionID: "dummy-function-id",
		}

		err := node.processRollCallResponse(context.Background(), mocks.GenericPeerID, res)
		require.NoError(t, err)

		// Verify roll call response is not found, even though the response has been processed.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-node.rollCall.responses(requestID):
			require.FailNow(t, "roll call response found but not expected")
		case <-ctx.Done():
			break
		}
	})
	t.Run("function install response", func(t *testing.T) {
		t.Parallel()

		msg := response.InstallFunction{
			Code:    codes.OK,
			Message: "dummy-message",
		}

		err := node.processInstallFunctionResponse(context.Background(), mocks.GenericPeerID, msg)
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

	t.Run("head node handles install", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.HeadNode)

		err := node.processInstallFunction(context.Background(), mocks.GenericPeerID, installReq)
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

			require.Equal(t, codes.Accepted, received.Code)
			require.Equal(t, expectedMessage, received.Message)
		})

		err = node.processInstallFunction(context.Background(), receiver.ID(), installReq)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("worker node handles function install error", func(t *testing.T) {
		t.Parallel()

		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)
		hostAddNewPeer(t, node.host, receiver)

		fstore := mocks.BaselineFStore(t)
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

		err = node.processInstallFunction(context.Background(), receiver.ID(), installReq)
		require.Error(t, err)
	})
	t.Run("worker node handles failure to send response", func(t *testing.T) {
		t.Parallel()

		// Receiver exists but not added to peer store - the node doesn't know
		// the receivers addresses so `send` will fail.
		receiver, err := host.New(mocks.NoopLogger, loopback, 0)
		require.NoError(t, err)

		node := createNode(t, blockless.WorkerNode)

		err = node.processInstallFunction(context.Background(), receiver.ID(), installReq)
		require.Error(t, err)
	})
}
