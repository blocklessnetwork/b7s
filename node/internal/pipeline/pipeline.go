package pipeline

import (
	"fmt"
)

type ID int

const (
	PubSub ID = iota + 1
	DirectMessage
)

func (i ID) String() string {
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
	ID    ID     // ID of the pipeline on which the message was received.
	Topic string // optional - topic on which this message was published.
}

func (p Pipeline) String() string {

	switch p.ID {
	case PubSub:
		return fmt.Sprintf("%v:%v", p.ID.String(), p.Topic)

	default:
		return p.ID.String()
	}
}

func DirectMessagePipeline() Pipeline {
	return Pipeline{ID: DirectMessage}
}

func PubSubPipeline(topic string) Pipeline {
	return Pipeline{PubSub, topic}
}
