package node

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node/internal/pipeline"
)

func TestNode_DisallowedMessages(t *testing.T) {

	var (
		pubsub = pipeline.PubSubPipeline(DefaultTopic)
		direct = pipeline.DirectMessagePipeline()
	)

	tests := []struct {
		pipeline pipeline.Pipeline
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
		err := allowedMessage(test.message, test.pipeline)
		require.ErrorIsf(t, err, errDisallowedMessage, "message: %s, pipeline: %s", test.message, test.pipeline)
	}
}
