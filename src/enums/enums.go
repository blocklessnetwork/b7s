package enums

var (
	MsgHealthCheck           = "MsgHealthCheck"
	MsgExecute               = "MsgExecute"
	MsgExecuteResult         = "MsgExecuteResult"
	MsgExecuteError          = "MsgExecuteError"
	MsgExecuteTimeout        = "MsgExecuteTimeout"
	MsgExecuteUnknown        = "MsgExecuteUnknown"
	MsgExecuteInvalid        = "MsgExecuteInvalid"
	MsgExecuteNotFound       = "MsgExecuteNotFound"
	MsgExecuteNotSupported   = "MsgExecuteNotSupported"
	MsgExecuteNotImplemented = "MsgExecuteNotImplemented"
	MsgExecuteNotAuthorized  = "MsgExecuteNotAuthorized"
	MsgExecuteNotPermitted   = "MsgExecuteNotPermitted"
	MsgExecuteNotAvailable   = "MsgExecuteNotAvailable"
	MsgExecuteNotReady       = "MsgExecuteNotReady"
	MsgExecuteNotConnected   = "MsgExecuteNotConnected"
	MsgExecuteNotInitialized = "MsgExecuteNotInitialized"
	MsgExecuteNotConfigured  = "MsgExecuteNotConfigured"
	MsgExecuteNotInstalled   = "MsgExecuteNotInstalled"
	MsgExecuteNotUpgraded    = "MsgExecuteNotUpgraded"
	MsgRollCall              = "MsgRollCall"
	MsgRollCallResponse      = "MsgRollCallResponse"
	MsgExecuteResponse       = "MsgExecuteResponse"
)

var (
	RequestExecute  = "RequestExecute"
	ResponseExecute = "ResponseExecute"
	RequestInstall  = "RequestInstall"
	ResponseInstall = "ResponseInstall"
)

var (
	ResponseCodeOk             = "200"
	ResponseCodeError          = "500"
	ResponseCodeTimeout        = "408"
	ResponseCodeUnknown        = "520"
	ResponseCodeInvalid        = "400"
	ResponseCodeNotFound       = "404"
	ResponseCodeNotSupported   = "501"
	ResponseCodeNotImplemented = "501"
	ResponseCodeNotAuthorized  = "401"
	ResponseCodeNotPermitted   = "403"
	ResponseCodeNotAvailable   = "503"
)
