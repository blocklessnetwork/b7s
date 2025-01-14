package head

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/request"
	"github.com/blessnetwork/b7s/models/response"
	"github.com/blessnetwork/b7s/testing/mocks"
)

func TestHead_Execute(t *testing.T) {

	var (
		requestID = fmt.Sprintf("request-id-%v", rand.Int())
		subgroup  = fmt.Sprintf("topic-%v", rand.Int())
		req       = mocks.GenericExecutionRequest
		res       = execute.NodeResult{
			Result: mocks.GenericExecutionResult,
		}

		workerID = mocks.GenericPeerID

		lock                 sync.Mutex
		start                time.Time
		rollCallPublished    time.Time
		executionRequestSent time.Time
	)

	head := createHeadNode(t)

	// We know the request ID so we'll setup a roll call reply in advance.
	head.rollCall.create(requestID)
	head.rollCall.add(
		requestID,
		rollCallResponse{
			From: workerID,
			RollCall: response.RollCall{
				Code:       codes.Accepted,
				FunctionID: req.FunctionID,
				RequestID:  requestID,
			},
		})

	core := mocks.BaselineNodeCore(t)
	core.ConnectedFunc = func(peer.ID) bool {
		return true
	}
	// Setup a publish func - we expect the head node to publish a roll call.
	core.PublishToTopicFunc = func(_ context.Context, topic string, msg blockless.Message) error {

		lock.Lock()
		defer lock.Unlock()

		require.Equal(t, subgroup, topic)

		rc, ok := any(msg).(*request.RollCall)
		require.True(t, ok)

		require.Equal(t, req.FunctionID, rc.FunctionID)
		require.Equal(t, requestID, rc.RequestID)

		c, err := consensus.Parse(req.Config.ConsensusAlgorithm)
		require.NoError(t, err)

		require.Equal(t, c, rc.Consensus)

		rollCallPublished = time.Now()

		return nil
	}

	// Setup a send func - we expect the head node to send the work order to the worker.
	core.SendToManyFunc = func(_ context.Context, peers []peer.ID, msg blockless.Message, requireAll bool) error {

		lock.Lock()
		defer lock.Unlock()

		require.Len(t, peers, 1)
		require.Equal(t, workerID, peers[0])
		require.False(t, requireAll) // Tailored for this test case - we're not using consensus so we don't requireAll.

		wo, ok := any(msg).(*request.WorkOrder)
		require.True(t, ok)

		require.Equal(t, requestID, wo.RequestID)
		require.Equal(t, req, wo.Request)

		executionRequestSent = time.Now()

		// Simulate getting a response after sending a work order.
		key := peerRequestKey(requestID, workerID)
		head.workOrderResponses.Set(key, res)

		return nil
	}

	head.Core = core

	// Main part of the test start time.
	start = time.Now()

	// Test the function execution.
	er := request.Execute{
		Request: req,
		Topic:   subgroup,
	}
	code, results, cluster, err := head.execute(context.Background(), requestID, er)
	require.NoError(t, err)

	// Verify code signals success.
	require.Equal(t, codes.OK, code)

	// Verify we have a single result and it's what we expect.
	require.Len(t, results, 1)
	require.Equal(t, res, results[workerID])

	// Verify we have a single node in the cluster and it's the worker we expect.
	require.Len(t, cluster.Peers, 1)
	require.Equal(t, workerID, cluster.Peers[0])

	// Verify actions happened in the order we expect them to.
	require.NotZero(t, start)
	require.True(t, rollCallPublished.After(start))
	require.True(t, executionRequestSent.After(rollCallPublished))
}

func createHeadNode(t *testing.T) *HeadNode {
	t.Helper()

	head, err := New(mocks.BaselineNodeCore(t))
	require.NoError(t, err)

	return head
}
