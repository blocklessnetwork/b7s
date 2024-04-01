package node

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func TestNode_DisallowedMessages(t *testing.T) {

	tests := []struct {
		message  string
		pipeline messagePipeline
	}{
		// Messages disallowed for publishing.
		{message: blockless.MessageInstallFunctionResponse, pipeline: subscriptionPipeline},
		{message: blockless.MessageExecute, pipeline: subscriptionPipeline},
		{message: blockless.MessageExecuteResponse, pipeline: subscriptionPipeline},
		{message: blockless.MessageFormCluster, pipeline: subscriptionPipeline},
		{message: blockless.MessageFormClusterResponse, pipeline: subscriptionPipeline},
		{message: blockless.MessageDisbandCluster, pipeline: subscriptionPipeline},

		// Messages disallowed for direct sending.
		{message: blockless.MessageHealthCheck, pipeline: directMessagePipeline},
		{message: blockless.MessageRollCall, pipeline: directMessagePipeline},
	}

	for _, test := range tests {
		err := allowedMessage(test.message, test.pipeline)
		require.ErrorIsf(t, err, errDisallowedMessage, "message: %s, pipeline: %s", test.message, test.pipeline)
	}
}
