package worker

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/request"
	"github.com/blessnetwork/b7s/models/response"
	"github.com/blessnetwork/b7s/testing/mocks"
)

func createWorkerNode(t *testing.T) *Worker {
	t.Helper()

	var (
		core     = mocks.BaselineNodeCore(t)
		executor = mocks.BaselineExecutor(t)
		fstore   = mocks.BaselineFStore(t)

		workspace = t.TempDir()
	)

	worker, err := New(core, fstore, executor, Workspace(workspace))
	require.NoError(t, err)

	return worker
}

func TestWorker_ProcessWorkOrder(t *testing.T) {

	// Create request and expected response.
	var (
		requestID = fmt.Sprintf("request-id-%v", rand.Int())
		req       = request.WorkOrder{
			RequestID: requestID,
			Request:   mocks.GenericExecutionRequest,
		}

		result = execute.Result{
			Code: codes.OK,
			Result: execute.RuntimeOutput{
				Stdout:   fmt.Sprintf("test-stdout-%v", rand.Int()),
				Stderr:   fmt.Sprintf("test-stderr-%v", rand.Int()),
				ExitCode: rand.Int(),
			},
			Usage: execute.Usage{
				WallClockTime: time.Duration(rand.Int()),
				CPUUserTime:   time.Duration(rand.Int()),
				CPUSysTime:    time.Duration(rand.Int()),
				MemoryMaxKB:   rand.Int64(),
			},
		}
	)

	// Create executor that verifies that input request is correct.
	executor := mocks.BaselineExecutor(t)
	executor.ExecFunctionFunc = func(ctx context.Context, id string, er execute.Request) (execute.Result, error) {
		require.Equal(t, id, requestID)
		require.Equal(t, req.Request, er)

		return result, nil
	}

	// Create node core with overrridden Send function.
	// Send function verifies that the result it is passed to it is what the executor returned.
	core := mocks.BaselineNodeCore(t)
	core.SendFunc = func(_ context.Context, _ peer.ID, msg bls.Message) error {
		er, ok := any(msg).(*response.WorkOrder)
		require.True(t, ok)

		require.Equal(t, result.Code, er.Result.Code)
		require.Equal(t, result.Result, er.Result.Result.Result) // RuntimeOutput
		require.Equal(t, result.Usage, er.Result.Result.Usage)

		return nil
	}

	worker := createWorkerNode(t)
	worker.executor = executor
	worker.Core = core

	err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
	require.NoError(t, err)
}

func TestWorker_ProcessWorkOrder_Metadata(t *testing.T) {

	var (
		req = request.WorkOrder{
			RequestID: "request-id",
			Request:   mocks.GenericExecutionRequest,
		}
	)
	data := make(map[string]any)
	keys := rand.IntN(10)
	for i := 0; i < keys; i++ {

		key := fmt.Sprintf("key-%v", rand.Int())

		// Switch up data - use both strings and ints.
		var value any
		if i%2 == 0 {
			value = fmt.Sprintf("value-%v", rand.Int())
		} else {
			value = rand.Int()
		}

		data[key] = value
	}

	// Create worker and set metadata provider.
	worker := createWorkerNode(t)
	mp := newDummyMetadataProvider(data)
	worker.cfg.MetadataProvider = mp

	// Setup Send function to verify metadata is correctly set
	core := mocks.BaselineNodeCore(t)
	core.SendFunc = func(_ context.Context, _ peer.ID, msg bls.Message) error {
		er, ok := any(msg).(*response.WorkOrder)
		require.True(t, ok)
		require.Equal(t, data, er.Result.Metadata)

		return nil
	}
	worker.Core = core

	err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
	require.NoError(t, err)

}

func TestWorker_ProcessWorkOrder_HandlesErrors(t *testing.T) {

	t.Run("function lookup error", func(t *testing.T) {

		var (
			req = request.WorkOrder{
				RequestID: "request-id",
				Request:   mocks.GenericExecutionRequest,
			}
		)

		worker := createWorkerNode(t)

		fstore := mocks.BaselineFStore(t)
		fstore.IsInstalledFunc = func(string) (bool, error) {
			return false, nil
		}
		worker.fstore = fstore

		// Override Send function to verify that the result it is passed to it is what the executor returned, and it was sent despite an execution error.
		core := mocks.BaselineNodeCore(t)
		core.SendFunc = func(_ context.Context, _ peer.ID, msg bls.Message) error {
			er, ok := any(msg).(*response.WorkOrder)
			require.True(t, ok)

			require.Equal(t, codes.NotFound, er.Result)
			require.Empty(t, er.Result.Result.Result) // RuntimeOutput
			require.Empty(t, er.Result.Result.Usage)

			return nil
		}

		err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
		require.NoError(t, err)
	})
	t.Run("request ID missing", func(t *testing.T) {

		var (
			req = request.WorkOrder{
				RequestID: "", // Empty
				Request:   mocks.GenericExecutionRequest,
			}
		)

		worker := createWorkerNode(t)

		err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
		require.Error(t, err)
	})
	t.Run("send error", func(t *testing.T) {

		var (
			sendErr = errors.New("send error")

			requestID = fmt.Sprintf("request-id-%v", rand.Int())
			req       = request.WorkOrder{
				RequestID: requestID,
				Request:   mocks.GenericExecutionRequest,
			}
		)

		worker := createWorkerNode(t)
		// Setup send failure.
		core := mocks.BaselineNodeCore(t)
		core.SendFunc = func(_ context.Context, _ peer.ID, _ bls.Message) error {
			return sendErr
		}
		worker.Core = core

		err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
		require.ErrorIs(t, err, sendErr)
	})
	t.Run("execution error", func(t *testing.T) {

		var (
			req = request.WorkOrder{
				RequestID: "request-id",
				Request:   mocks.GenericExecutionRequest,
			}

			result = execute.Result{
				Code: codes.OK,
				Result: execute.RuntimeOutput{
					Stdout:   fmt.Sprintf("test-stdout-%v", rand.Int()),
					Stderr:   fmt.Sprintf("test-stderr-%v", rand.Int()),
					ExitCode: rand.Int(),
				},
				Usage: execute.Usage{
					WallClockTime: time.Duration(rand.Int()),
					CPUUserTime:   time.Duration(rand.Int()),
					CPUSysTime:    time.Duration(rand.Int()),
					MemoryMaxKB:   rand.Int64(),
				},
			}
		)

		worker := createWorkerNode(t)

		// Create an executor where execution fails.
		executor := mocks.BaselineExecutor(t)
		executor.ExecFunctionFunc = func(_ context.Context, _ string, _ execute.Request) (execute.Result, error) {
			return result, errors.New("execution error")
		}
		worker.executor = executor

		// Override Send function to verify that the result it is passed to it is what the executor returned, and it was sent despite an execution error.
		core := mocks.BaselineNodeCore(t)
		core.SendFunc = func(_ context.Context, _ peer.ID, msg bls.Message) error {
			er, ok := any(msg).(*response.WorkOrder)
			require.True(t, ok)

			require.Equal(t, result.Code, er.Result.Code)
			require.Equal(t, result.Result, er.Result.Result.Result) // RuntimeOutput
			require.Equal(t, result.Usage, er.Result.Result.Usage)

			return nil
		}

		err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
		require.NoError(t, err)
	})
	t.Run("consensus required but no cluster", func(t *testing.T) {

		var (
			req = request.WorkOrder{
				RequestID: "request-id",
				Request:   mocks.GenericExecutionRequest,
			}
		)

		req.Config.ConsensusAlgorithm = consensus.Raft.String()

		worker := createWorkerNode(t)

		core := mocks.BaselineNodeCore(t)
		core.SendFunc = func(_ context.Context, _ peer.ID, msg bls.Message) error {
			er, ok := any(msg).(*response.WorkOrder)
			require.True(t, ok)

			require.Equal(t, codes.Error, er.Code)
			require.Empty(t, er.Result.Result.Result) // RuntimeOutput
			require.Empty(t, er.Result.Result.Usage)

			return nil
		}
		worker.Core = core

		err := worker.processWorkOrder(context.Background(), mocks.GenericPeerID, req)
		require.NoError(t, err)
	})
}

type dummyMetadataProvider struct {
	data any
}

func newDummyMetadataProvider(data any) *dummyMetadataProvider {

	return &dummyMetadataProvider{data: data}
}

func (m *dummyMetadataProvider) Metadata(_ execute.Request, _ execute.RuntimeOutput) (any, error) {
	return m.data, nil
}
