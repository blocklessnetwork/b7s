package enums

import (
	"github.com/libp2p/go-libp2p/core/protocol"
)

var (
	MsgHealthCheck             = "MsgHealthCheck"
	MsgExecute                 = "MsgExecute"
	MsgExecuteResult           = "MsgExecuteResult"
	MsgExecuteError            = "MsgExecuteError"
	MsgExecuteTimeout          = "MsgExecuteTimeout"
	MsgExecuteUnknown          = "MsgExecuteUnknown"
	MsgExecuteInvalid          = "MsgExecuteInvalid"
	MsgExecuteNotFound         = "MsgExecuteNotFound"
	MsgExecuteNotSupported     = "MsgExecuteNotSupported"
	MsgExecuteNotImplemented   = "MsgExecuteNotImplemented"
	MsgExecuteNotAuthorized    = "MsgExecuteNotAuthorized"
	MsgExecuteNotPermitted     = "MsgExecuteNotPermitted"
	MsgExecuteNotAvailable     = "MsgExecuteNotAvailable"
	MsgExecuteNotReady         = "MsgExecuteNotReady"
	MsgExecuteNotConnected     = "MsgExecuteNotConnected"
	MsgExecuteNotInitialized   = "MsgExecuteNotInitialized"
	MsgExecuteNotConfigured    = "MsgExecuteNotConfigured"
	MsgExecuteNotInstalled     = "MsgExecuteNotInstalled"
	MsgExecuteNotUpgraded      = "MsgExecuteNotUpgraded"
	MsgRollCall                = "MsgRollCall"
	MsgRollCallResponse        = "MsgRollCallResponse"
	MsgExecuteResponse         = "MsgExecuteResponse"
	MsgInstallFunction         = "MsgInstallFunction"
	MsgInstallFunctionResponse = "MsgInstallFunctionResponse"
)

var (
	RequestExecute  = "RequestExecute"
	ResponseExecute = "ResponseExecute"
	RequestInstall  = "RequestInstall"
	ResponseInstall = "ResponseInstall"
)

var (
	ResponseCodeOk             = "200"
	ResponseCodeAccepted       = "202"
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

var (
	WorkProtocolId protocol.ID = "/b7s/work/1.0.0"
)

var (
	RoleWorker = "worker"
	RoleHead   = "head"
)

var (
	ChannelMsgLocal = "ChannelMsgLocal"
	// ChannelMsgInstallFunction  = "ChannelMsgInstallFunction"
	// ChannelMsgExecute          = "ChannelMsgExecute"
	ChannelMsgExecuteResponse = "ChannelMsgExecuteResponse"
	// ChannelMsgRollCall         = "ChannelMsgRollCall"
	// ChannelMsgHealthCheck      = "ChannelMsgHealthCheck"
	ChannelMsgRollCallResponse = "ChannelMsgRollCallResponse"
)
