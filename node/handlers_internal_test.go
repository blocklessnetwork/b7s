package node

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/testing/mocks"
	"github.com/stretchr/testify/require"
)

func TestNode_Handlers(t *testing.T) {

	node := createNode(t, blockless.HeadNode)

	t.Run("health check", func(t *testing.T) {
		t.Parallel()

		msg := response.Health{
			Type: blockless.MessageHealthCheck,
			Code: response.CodeOK,
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
			Code:       response.CodeOK,
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
			Code:    response.CodeOK,
			Message: "dummy-message",
		}

		payload := serialize(t, msg)
		err := node.processInstallFunctionResponse(context.Background(), mocks.GenericPeerID, payload)
		require.NoError(t, err)
	})
}

func serialize(t *testing.T, message any) []byte {
	payload, err := json.Marshal(message)
	require.NoError(t, err)

	return payload
}
