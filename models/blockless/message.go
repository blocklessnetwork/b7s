package blockless

import (
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

type Message interface {
	Type() string
}

// Message types in the Blockless protocol.
const (
	MessageHealthCheck             = "MsgHealthCheck"
	MessageInstallFunction         = "MsgInstallFunction"
	MessageInstallFunctionResponse = "MsgInstallFunctionResponse"
	MessageRollCall                = "MsgRollCall"
	MessageRollCallResponse        = "MsgRollCallResponse"
	MessageExecute                 = "MsgExecute"
	MessageExecuteResponse         = "MsgExecuteResponse"
	MessageFormCluster             = "MsgFormCluster"
	MessageFormClusterResponse     = "MsgFormClusterResponse"
	MessageDisbandCluster          = "MsgDisbandCluster"
)

type TraceableMessage interface {
	Message
	SaveTraceContext(tracing.TraceInfo)
}

type BaseMessage struct {
	tracing.TraceInfo
}

func (m *BaseMessage) SaveTraceContext(t tracing.TraceInfo) {
	m.TraceInfo = t
}
