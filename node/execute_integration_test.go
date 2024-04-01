//go:build integration
// +build integration

package node_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

func TestNode_ExecuteComplete(t *testing.T) {

	const (
		testTimeLimit = 1 * time.Minute

		dirPattern = "b7s-node-execute-integration-test-"

		cid = "whatever-cid"

		// Paths where files will be hosted on the test server.
		manifestEndpoint    = "/hello-manifest.json"
		archiveEndpoint     = "/hello-deployment.tar.gz"
		testFunctionToServe = "testdata/hello.tar.gz"
		functionMethod      = "hello.wasm"

		expectedExecutionResult = `This is the start of my program
The answer is  42
This is the end of my program
`
	)

	cleanupDisabled := cleanupDisabled()

	var verifiedExecution bool

	t.Log("starting test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set a hard limit for test duration.
	// This looks a bit sketchy as tests can have the time limit
	// set externally, but as there's a lot of moving pieces here,
	// include it for better usability.
	go func() {
		<-time.After(testTimeLimit)
		cancel()
		t.Log("cancelling test")
	}()

	// Phase 0: Create libp2p hosts, loggers, temporary directories and nodes.

	head := instantiateNode(t, dirPattern, blockless.HeadNode)
	defer head.db.Close()
	defer head.logFile.Close()
	if !cleanupDisabled {
		defer os.RemoveAll(head.dir)
	}

	worker := instantiateNode(t, dirPattern, blockless.WorkerNode)
	defer worker.db.Close()
	defer worker.logFile.Close()
	if !cleanupDisabled {
		defer os.RemoveAll(worker.dir)
	}

	t.Log("created nodes")

	// Phase 1: Setup connections and start node main loops.

	// Client that will issue and receive request.
	client := createClient(t)

	// Add hosts to each others peer stores so that they know how to contact each other.
	hostAddNewPeer(t, client.host, head.host)
	hostAddNewPeer(t, client.host, worker.host)
	hostAddNewPeer(t, head.host, worker.host)

	// Establish a connection so that hosts disseminate topic subscription info.
	headInfo := hostGetAddrInfo(t, head.host)
	err := worker.host.Connect(ctx, *headInfo)
	require.NoError(t, err)

	t.Log("setup addressing")

	// Phase 2: Start nodes.

	t.Log("starting nodes")

	// We start nodes in separate goroutines.
	var nodesWG sync.WaitGroup
	nodesWG.Add(1)
	go func() {
		defer nodesWG.Done()

		err := head.node.Run(ctx)
		require.NoError(t, err)

		t.Log("head node stopped")
	}()
	nodesWG.Add(1)
	go func() {
		defer nodesWG.Done()

		err := worker.node.Run(ctx)
		require.NoError(t, err)

		t.Log("worker node stopped")
	}()

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
	err = client.sendInstallMessage(ctx, worker.host.ID(), manifestURL, cid)
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
		require.Equal(t, expectedExecutionResult, res.Results[worker.host.ID()].Result.Stdout)

		t.Log("client verified execution response")

		verifiedExecution = true
	})

	err = client.sendExecutionMessage(ctx, head.host.ID(), cid, functionMethod, 0, 1)
	require.NoError(t, err)

	executeWG.Wait()

	t.Log("execution request processed")

	// Since we're done, we can cancel the context, leading to stopping of the nodes.
	cancel()

	nodesWG.Wait()

	t.Log("nodes shutdown")

	t.Log("test complete")

	require.True(t, verifiedExecution)
}
