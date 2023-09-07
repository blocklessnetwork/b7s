//go:build integration
// +build integration

package node_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/consensus"
	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/response"
)

func TestNode_PBFT_ExecuteComplete(t *testing.T) {

	const (
		testTimeLimit = 1 * time.Minute

		dirPattern = "b7s-node-pbft-execute-integration-test-"

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
	nodeDir := fmt.Sprintf("%v-head-", dirPattern)
	head := instantiateNode(t, nodeDir, blockless.HeadNode)
	t.Logf("head node workspace: %s", head.dir)

	var workers []*nodeScaffolding
	for i := 0; i < 4; i++ {
		nodeDir := fmt.Sprintf("%v-worker-%v-", dirPattern, i)

		worker := instantiateNode(t, nodeDir, blockless.WorkerNode)
		t.Logf("worker node #%v workspace: %s", i, worker.dir)

		workers = append(workers, worker)
	}

	workerIDs := make([]peer.ID, 0, len(workers))
	for _, worker := range workers {
		workerIDs = append(workerIDs, worker.host.ID())
	}

	// Cleanup everything after test is complete.
	defer func() {
		head.db.Close()
		head.logFile.Close()
		if !cleanupDisabled {
			os.RemoveAll(head.dir)
		}

		for _, worker := range workers {
			worker.db.Close()
			worker.logFile.Close()
			if !cleanupDisabled {
				os.RemoveAll(worker.dir)
			}
		}
	}()

	var nodes []*nodeScaffolding
	nodes = append(nodes, head)
	nodes = append(nodes, workers...)

	t.Log("created nodes")

	// Phase 1: Setup connections.

	// Client that will issue and receive request.
	client := createClient(t)

	// Add hosts to each others peer stores so that they know how to contact each other, and then establish connections.
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if j == i {
				continue
			}
			hostAddNewPeer(t, client.host, nodes[i].host)
			hostAddNewPeer(t, nodes[i].host, nodes[j].host)
			hostAddNewPeer(t, nodes[j].host, nodes[i].host)

			// Establish a connection so that hosts disseminate topic subscription info.
			info := hostGetAddrInfo(t, nodes[j].host)
			err := nodes[i].host.Connect(ctx, *info)
			require.NoError(t, err)
		}
	}

	// Phase 2: Start nodes.
	t.Log("starting nodes")

	// We start nodes in separate goroutines.
	var nodesWG sync.WaitGroup
	nodesWG.Add(len(nodes))

	for _, node := range nodes {
		go func(node *nodeScaffolding) {
			defer nodesWG.Done()

			err := node.node.Run(ctx)
			require.NoError(t, err)

			t.Log("node stopped")
		}(node)
	}

	// Add a delay for the hosts to subscribe to topics,
	// diseminate subscription information etc.
	time.Sleep(startupDelay)

	t.Log("starting function server")

	// Phase 3: Create the server hosting the manifest and the function.

	srv := createFunctionServer(t, manifestEndpoint, archiveEndpoint, testFunctionToServe, cid)
	defer srv.Close()

	// Phase 4: Have the worker nodes install the function.
	// That way, when he receives the execution request - he will have the function needed to execute it.

	t.Log("instructing worker node to install function")

	var installWG sync.WaitGroup
	installWG.Add(len(workers))

	// Setup verifier for the response we expect.
	client.host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer installWG.Done()
		defer stream.Close()

		from := stream.Conn().RemotePeer()
		require.Contains(t, workerIDs, from)

		var res response.InstallFunction
		getStreamPayload(t, stream, &res)

		require.Equal(t, blockless.MessageInstallFunctionResponse, res.Type)
		require.Equal(t, codes.Accepted, res.Code)
		require.Equal(t, "installed", res.Message)

		t.Log("client received function install response")
	})

	manifestURL := fmt.Sprintf("%v%v", srv.URL, manifestEndpoint)
	for _, worker := range workers {
		err := client.sendInstallMessage(ctx, worker.host.ID(), manifestURL, cid)
		require.NoError(t, err)
	}

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

		payload, _ := json.Marshal(res)
		fmt.Printf("%s\n", payload)

		require.Equal(t, blockless.MessageExecuteResponse, res.Type)
		require.Equal(t, codes.OK, res.Code)
		require.NotEmpty(t, res.RequestID)

		require.Len(t, res.Cluster.Peers, len(workers))

		// Verify cluster nodes are the workers we created.
		require.ElementsMatch(t, workerIDs, res.Cluster.Peers)

		require.GreaterOrEqual(t, uint(len(res.Results)), pbft.MinClusterResults(uint(len(workers))))

		for peer, exres := range res.Results {
			require.Contains(t, workerIDs, peer)
			require.Equal(t, expectedExecutionResult, exres.Result.Stdout)
		}

		t.Log("client verified execution response")

		verifiedExecution = true
	})

	err := client.sendExecutionMessage(ctx, head.host.ID(), cid, functionMethod, consensus.PBFT, len(workers))
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
