package node

import (
	"github.com/blocklessnetwork/b7s/models/blockless"
	pp "github.com/blocklessnetwork/b7s/node/internal/pipeline"
)

func messageAllowedOnPipeline(msg string, pipeline pp.Pipeline) bool {

	if pipeline.ID == pp.DirectMessage {

		switch msg {
		// Messages we don't expect as direct messages.
		case
			blockless.MessageHealthCheck,
			blockless.MessageRollCall:

			// Technically we only publish InstallFunction. However, it's handy for tests to support
			// direct install, and it's somewhat of a low risk.

			return false

		default:
			return true
		}
	}

	switch msg {
	// Messages we don't allow to be published.
	case
		blockless.MessageInstallFunctionResponse,
		blockless.MessageExecute,
		blockless.MessageExecuteResponse,
		blockless.MessageFormCluster,
		blockless.MessageFormClusterResponse,
		blockless.MessageDisbandCluster,
		blockless.MessageRollCallResponse:

		return false

	default:
		return true
	}
}
