package blockless

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

type Message interface {
	Type() string
}
