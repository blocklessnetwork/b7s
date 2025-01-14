//go:build integration
// +build integration

package node

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/response"
	"github.com/blessnetwork/b7s/testing/helpers"
)

func TestHeadNode_Execute(t *testing.T) {

	const (
		dirPattern    = "b7s-node-execute-integration-test-"
		testTimeLimit = 1 * time.Minute

		// Paths where files will be hosted on the test server.
		manifestEndpoint    = "/hello-manifest.json"
		archiveEndpoint     = "/hello-deployment.tar.gz"
		testFunctionToServe = "testdata/hello.tar.gz"
		cid                 = "whatever-cid"
		functionMethod      = "hello.wasm"

		expectedExecutionResult = `This is the start of my program
The answer is  42
This is the end of my program
`
	)

	var (
		cleanupDisabled = cleanupDisabled()
		ctx, cancel     = context.WithTimeout(context.Background(), testTimeLimit)

		verifiedExecution atomic.Bool
	)
	defer cancel()

	t.Log("starting test")

	// Phase 0: Create libp2p hosts, loggers, temporary directories and nodes.

	headNode := instantiateNode(t, dirPattern, blockless.HeadNode)
	defer headNode.logFile.Close()
	if !cleanupDisabled {
		defer os.RemoveAll(headNode.dir)
	}

	workerNode := instantiateNode(t, dirPattern, blockless.WorkerNode)
	defer workerNode.db.Close()
	defer workerNode.logFile.Close()
	if !cleanupDisabled {
		defer os.RemoveAll(workerNode.dir)
	}

	t.Log("created nodes")

	// Phase 1: Setup connections and start node main loops.

	client := createClient(t)

	helpers.HostAddNewPeer(t, client.host, headNode.host)
	helpers.HostAddNewPeer(t, client.host, workerNode.host)
	helpers.HostAddNewPeer(t, headNode.host, workerNode.host)

	// Establish a connection so that hosts disseminate topic subscription info.
	err := workerNode.host.Connect(ctx, *helpers.HostGetAddrInfo(t, headNode.host))
	require.NoError(t, err)

	t.Log("setup addressing")

	// Phase 2: Start nodes.

	t.Log("starting nodes")

	var runErr multierror.Group
	// We require Run to not fail below so that we can scrap a test earlier if something goes wrong.
	runErr.Go(func() error {
		err := headNode.node.Run(ctx)
		require.NoError(t, err)
		return err

	})
	runErr.Go(func() error {
		err := workerNode.node.Run(ctx)
		require.NoError(t, err)
		return err
	})

	// Add a delay for the hosts to subscribe to topics,
	// diseminate subscription information etc.
	time.Sleep(startupDelay)

	t.Log("starting function server")

	// Phase 3: Create the server hosting the manifest and the function.

	srv := createFunctionServer(t, manifestEndpoint, archiveEndpoint, testFunctionToServe, cid)
	defer srv.Close()

	// Phase 4: Have the worker install the function.
	// That way, when he receives the execution request - he will have the function needed to execute it.

	t.Log("instructing worker node to install function")

	var installWG sync.WaitGroup
	installWG.Add(1)

	// Setup verifier for the response we expect.
	client.host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer installWG.Done()
		defer stream.Close()

		var res response.InstallFunction
		getStreamPayload(t, stream, &res)

		require.Equal(t, codes.Accepted, res.Code)
		require.Equal(t, "installed", res.Message)

		t.Log("client received function install response")
	})

	manifestURL := fmt.Sprintf("%v%v", srv.URL, manifestEndpoint)
	err = client.sendInstallMessage(ctx, workerNode.host.ID(), manifestURL, cid)
	require.NoError(t, err)

	// Wait for the installation request to be processed.
	installWG.Wait()

	t.Log("worker node installed function")

	// Phase 5: Request execution from the head node.

	t.Log("sending execution request")

	// Setup verifier for the response we expect.
	var executeWG sync.WaitGroup

	executeWG.Add(1)
	client.host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer executeWG.Done()
		defer stream.Close()

		t.Log("client received execution response")

		var res response.Execute
		getStreamPayload(t, stream, &res)

		require.Equal(t, codes.OK, res.Code)
		require.NotEmpty(t, res.RequestID)
		require.Equal(t, expectedExecutionResult, res.Results[workerNode.host.ID()].Result.Result.Stdout)

		t.Log("client verified execution response")

		verifiedExecution.Store(true)
	})

	err = client.sendExecutionMessage(ctx, headNode.host.ID(), cid, functionMethod, 0, 1)
	require.NoError(t, err)

	executeWG.Wait()

	t.Log("execution request processed")

	// Since we're done, we can cancel the context, leading to stopping of the nodes.
	cancel()

	err = runErr.Wait().ErrorOrNil()
	require.NoError(t, err)

	t.Log("nodes shutdown")

	require.True(t, verifiedExecution.Load())

	t.Log("test complete")
}
