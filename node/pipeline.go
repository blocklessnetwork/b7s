package node

import (
	"errors"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type messagePipeline int

const (
	subscriptionPipeline messagePipeline = iota + 1
	directMessagePipeline
)

var errDisallowedMessage = errors.New("disallowed message")

func (p messagePipeline) String() string {
	switch p {
	case subscriptionPipeline:
		return "Subscription"
	case directMessagePipeline:
		return "DirectMessage"
	default:
		return "Unknown"
	}
}

func allowedMessage(msg string, pipeline messagePipeline) error {

	if pipeline == directMessagePipeline {

		switch msg {
		// Messages we don't expect as direct messages.
		case
			blockless.MessageHealthCheck,
			blockless.MessageRollCall:

			// Technically we only publish InstallFunction. However, it's handy for tests to support
			// direct install, and it's somewhat of a low risk.

			return errDisallowedMessage

		default:
			return nil
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

		return errDisallowedMessage

	default:
		return nil
	}
}
