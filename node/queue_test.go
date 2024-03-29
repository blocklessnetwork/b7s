package node

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/response"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestRollCallQueue(t *testing.T) {

	var (
		requestID = "dummy-request-id"

		res = rollCallResponse{
			From: mocks.GenericPeerID,
			RollCall: response.RollCall{
				RequestID:  requestID,
				FunctionID: "dummy-function-id",
			},
		}
	)

	t.Run("roll call queue works", func(t *testing.T) {

		const (
			count = 20
		)

		queue := newQueue(100)

		// Request does not exist in an empty map.
		require.False(t, queue.exists(requestID))

		// Request exists after creation.
		queue.create(requestID)
		require.True(t, queue.exists(requestID))

		var wg sync.WaitGroup
		wg.Add(count)
		//  Add `count` responses in parallel.
		for i := 0; i < count; i++ {
			go func() {
				defer wg.Done()
				queue.add(requestID, res)
			}()
		}
		wg.Wait()

		// Verify we have all responses recorded.
		responses := queue.responses(requestID)
		require.Len(t, responses, count)

		for i := 0; i < count; i++ {
			r := <-responses
			require.Equal(t, res, r)
		}

		queue.remove(requestID)
		require.False(t, queue.exists(requestID))
	})
	t.Run("roll call with pending responses removal works", func(t *testing.T) {

		const (
			count = 5
		)

		queue := newQueue(100)
		queue.create(requestID)

		for i := 0; i < count; i++ {
			queue.add(requestID, res)
		}

		queue.remove(requestID)
		require.False(t, queue.exists(requestID))
	})
}
