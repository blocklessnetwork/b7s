package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_RestExecute(t *testing.T) {

	var (
		request = mocks.GenericExecutionRequest

		result = execute.Result{
			Code: codes.OK,
			Result: execute.RuntimeOutput{
				Stdout:   "executor-output",
				Stderr:   "executor stderr log",
				ExitCode: 101,
			},
		}
	)

	node := createNode(t, blockless.WorkerNode)

	executor := mocks.BaselineExecutor(t)
	executor.ExecFunctionFunc = func(requestID string, req execute.Request) (execute.Result, error) {
		return result, nil
	}
	node.executor = executor

	code, results, err := node.ExecuteFunction(context.Background(), request)
	require.NoError(t, err)

	require.Equal(t, result.Code, code)
	require.Len(t, results, 1)

	// Worker node should provide only its own result.
	nodeID := node.host.ID().String()
	res, ok := results[nodeID]
	require.True(t, ok)
	require.Equal(t, result, res)
}

func TestNode_InstallMessageFromCID(t *testing.T) {

	const (
		cid                 = "dummy-cid"
		expectedManifestURL = "https://dummy-cid.ipfs.w3s.link/manifest.json"
	)

	req := createInstallMessageFromCID(cid)

	require.Equal(t, blockless.MessageInstallFunction, req.Type)
	require.Equal(t, cid, req.CID)
	require.Equal(t, expectedManifestURL, req.ManifestURL)
}
