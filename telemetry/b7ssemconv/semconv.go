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

const (
	FunctionCID    = attribute.Key("function.cid")
	FunctionMethod = attribute.Key("function.method")
)

const (
	ExecutionNodeCount = attribute.Key("execution.node.count")
	ExecutionConsensus = attribute.Key("execution.consensus")
	ExecutionRequestID = attribute.Key("execution.request.id")
)
