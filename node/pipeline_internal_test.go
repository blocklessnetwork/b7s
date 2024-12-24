package node

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func TestNode_DisallowedMessages(t *testing.T) {

	var (
		pubsub = PubSubPipeline(blockless.DefaultTopic)
		direct = DirectMessagePipeline
	)

	tests := []struct {
		pipeline Pipeline
		message  string
	}{
		// Messages disallowed for publishing.
		{pubsub, blockless.MessageInstallFunctionResponse},
		{pubsub, blockless.MessageExecute},
		{pubsub, blockless.MessageExecuteResponse},
		{pubsub, blockless.MessageFormCluster},
		{pubsub, blockless.MessageFormClusterResponse},
		{pubsub, blockless.MessageDisbandCluster},
		// Messages disallowed for direct sending.
		{direct, blockless.MessageHealthCheck},
		{direct, blockless.MessageRollCall},
	}

	for _, test := range tests {
		ok := correctPipeline(test.message, test.pipeline)
		require.False(t, ok, "message: %s, pipeline: %s", test.message, test.pipeline)
	}
}
