package b7ssemconv

import (
	"go.opentelemetry.io/otel/attribute"
)

// TODO: Add documentation for these.

const (
	MessageType     = attribute.Key("message.type")
	MessageID       = attribute.Key("message.id")
	MessagePipeline = attribute.Key("message.pipeline")
	MessagePeer     = attribute.Key("message.peer")
)
