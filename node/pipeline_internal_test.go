package node

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/models/bls"
)

func TestNode_DisallowedMessages(t *testing.T) {

	var (
		pubsub = PubSubPipeline(bls.DefaultTopic)
		direct = DirectMessagePipeline
	)

	tests := []struct {
		pipeline Pipeline
		message  string
	}{
		// Messages disallowed for publishing.
		{pubsub, bls.MessageInstallFunctionResponse},
		{pubsub, bls.MessageExecute},
		{pubsub, bls.MessageExecuteResponse},
		{pubsub, bls.MessageFormCluster},
		{pubsub, bls.MessageFormClusterResponse},
		{pubsub, bls.MessageDisbandCluster},
		// Messages disallowed for direct sending.
		{direct, bls.MessageHealthCheck},
		{direct, bls.MessageRollCall},
	}

	for _, test := range tests {
		ok := correctPipeline(test.message, test.pipeline)
		require.False(t, ok, "message: %s, pipeline: %s", test.message, test.pipeline)
	}
}
