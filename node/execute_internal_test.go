package node

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

// TODO: Responses should not have a "from" field

func TestNode_WorkerExecute(t *testing.T) {

	const (
		address = "127.0.0.1"
		port    = 0

		functionID     = "dummy-function-id"
		functionMethod = "dummy-function-method"
	)

	executionRequest := request.Execute{
		Type:       blockless.MessageExecute,
		FunctionID: functionID,
		Method:     functionMethod,
		Parameters: []execute.Parameter{},
		Config:     execute.Config{},
	}

	payload, err := json.Marshal(executionRequest)
	require.NoError(t, err)

	t.Run("handles correct execution", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		// Use a custom executor to verify all execution parameters are correct.
		executor := mocks.BaselineExecutor(t)
		executor.ExecFunctionFunc = func(req execute.Request) (execute.Result, error) {
			require.Equal(t, executionRequest.FunctionID, req.FunctionID)
			require.Equal(t, executionRequest.Method, req.Method)
			require.ElementsMatch(t, executionRequest.Parameters, req.Parameters)
			require.Equal(t, executionRequest.Config, req.Config)

			return mocks.GenericExecutionResult, nil
		}
		node.execute = executor

		// Create a host that will serve as a receiver of the execution response.
		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		recvAddr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, recvAddr)

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received response.Execute
			getStreamPayload(t, stream, &received)

			require.Equal(t, blockless.MessageExecuteResponse, received.Type)

			// We should receive the response the baseline executor will return.
			expected := mocks.GenericExecutionResult
			require.Equal(t, expected.RequestID, received.RequestID)
			require.Equal(t, expected.Code, received.Code)
			require.Equal(t, expected.Result, received.Result)
		})

		err = node.processExecute(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("handles execution failure", func(t *testing.T) {
		t.Parallel()

		var (
			faultyExecutionResult = execute.Result{
				Code:      response.CodeError,
				Result:    "something horrible has happened",
				RequestID: mocks.GenericUUID.String(),
			}
		)

		node := createNode(t, blockless.WorkerNode)

		// Use a custom executor to verify all execution parameters are correct.
		executor := mocks.BaselineExecutor(t)
		executor.ExecFunctionFunc = func(req execute.Request) (execute.Result, error) {
			return faultyExecutionResult, mocks.GenericError
		}
		node.execute = executor

		// Create a host that will serve as a receiver of the execution response.
		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		recvAddr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, recvAddr)

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received response.Execute
			getStreamPayload(t, stream, &received)

			require.Equal(t, blockless.MessageExecuteResponse, received.Type)

			require.Equal(t, faultyExecutionResult.RequestID, received.RequestID)
			require.Equal(t, faultyExecutionResult.Code, received.Code)
			require.Equal(t, faultyExecutionResult.Result, received.Result)
		})

		err = node.processExecute(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})
	t.Run("handles function store errors", func(t *testing.T) {
		t.Parallel()

		node := createNode(t, blockless.WorkerNode)

		// Error retrieving function.
		fstore := mocks.BaselineFunctionHandler(t)
		fstore.GetFunc = func(string, string, bool) (*blockless.FunctionManifest, error) {
			return nil, mocks.GenericError
		}
		node.function = fstore

		// Create a host that will serve as a receiver of the execution response.
		receiver, err := host.New(mocks.NoopLogger, address, port)
		require.NoError(t, err)

		recvAddr := getHostAddr(t, receiver)
		addPeerToPeerStore(t, node.host, recvAddr)

		var wg sync.WaitGroup
		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received response.Execute
			getStreamPayload(t, stream, &received)

			require.Equal(t, blockless.MessageExecuteResponse, received.Type)

			require.Equal(t, received.Code, response.CodeError)
		})

		err = node.processExecute(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()

		// Function is not installed.
		fstore.GetFunc = func(string, string, bool) (*blockless.FunctionManifest, error) {
			return nil, blockless.ErrNotFound
		}
		node.function = fstore

		wg.Add(1)

		receiver.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
			defer wg.Done()
			defer stream.Close()

			from := stream.Conn().RemotePeer()
			require.Equal(t, node.host.ID(), from)

			var received response.Execute
			getStreamPayload(t, stream, &received)

			require.Equal(t, blockless.MessageExecuteResponse, received.Type)

			require.Equal(t, received.Code, response.CodeNotFound)
		})

		err = node.processExecute(context.Background(), receiver.ID(), payload)
		require.NoError(t, err)

		wg.Wait()
	})

}
