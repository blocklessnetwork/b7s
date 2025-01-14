package node

import (
	"fmt"

	"github.com/blessnetwork/b7s/models/bls"
)

type PipelineID int

const (
	PubSub PipelineID = iota + 1
	DirectMessage
)

func (i PipelineID) String() string {
	switch i {
	case PubSub:
		return "pubsub"
	case DirectMessage:
		return "direct"
	default:
		return "unknown"
	}
}

type Pipeline struct {
	ID    PipelineID // ID of the pipeline on which the message was received.
	Topic string     // optional - topic on which this message was published.
}

var DirectMessagePipeline = Pipeline{ID: DirectMessage}

func (p Pipeline) String() string {

	switch p.ID {
	case PubSub:
		return fmt.Sprintf("%v:%v", p.ID.String(), p.Topic)

	default:
		return p.ID.String()
	}
}

func PubSubPipeline(topic string) Pipeline {
	return Pipeline{PubSub, topic}
}

func correctPipeline(msg string, pipeline Pipeline) bool {

	if pipeline.ID == DirectMessage {

		switch msg {
		// Messages we don't expect as direct messages.
		case
			bls.MessageHealthCheck,
			bls.MessageRollCall:

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
		bls.MessageInstallFunctionResponse,
		bls.MessageExecute,
		bls.MessageExecuteResponse,
		bls.MessageFormCluster,
		bls.MessageFormClusterResponse,
		bls.MessageDisbandCluster,
		bls.MessageRollCallResponse:

		return false

	default:
		return true
	}
}
