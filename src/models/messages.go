package models

import (
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Message struct {
	Type string
	Data interface{}
}

type MsgBase struct {
	Type string  `json:"type,omitempty"`
	From peer.ID `json:"from,omitempty"`
}

type MsgHealthPing struct {
	Type string  `json:"type,omitempty"`
	From peer.ID `json:"from,omitempty"`
	Code string  `json:"code,omitempty"`
}

func NewMsgHealthPing(code string) *MsgHealthPing {
	return &MsgHealthPing{
		Type: enums.MsgHealthCheck,
		Code: code,
	}
}

type MsgExecute struct {
	Type       string                     `json:"type,omitempty"`
	From       peer.ID                    `json:"from,omitempty"`
	Code       string                     `json:"code,omitempty"`
	FunctionId string                     `json:"functionId,omitempty"`
	Method     string                     `json:"method,omitempty"`
	Parameters []RequestExecuteParameters `json:"parameters,omitempty"`
	Config     ExecutionRequestConfig     `json:"config,omitempty"`
}

func NewMsgExecute(code string) *MsgExecute {
	return &MsgExecute{
		Type: enums.MsgExecute,
		Code: code,
	}
}

type MsgExecuteResponse struct {
	Type      string  `json:"type,omitempty"`
	RequestId string  `json:"requestId,omitempty"`
	From      peer.ID `json:"from,omitempty"`
	Code      string  `json:"code,omitempty"`
	Result    string  `json:"result,omitempty"`
}

type MsgRollCall struct {
	From       peer.ID `json:"from,omitempty"`
	Type       string  `json:"type,omitempty"`
	FunctionId string  `json:"functionId,omitempty"`
	RequestId  string  `json:"request_id,omitempty"`
}

func NewMsgRollCall(functionId string) *MsgRollCall {
	requestId, _ := uuid.NewRandom()
	return &MsgRollCall{
		Type:       enums.MsgRollCall,
		FunctionId: functionId,
		RequestId:  requestId.String(),
	}
}

type MsgRollCallResponse struct {
	Type       string  `json:"type,omitempty"`
	From       peer.ID `json:"from,omitempty"`
	Code       string  `json:"code,omitempty"`
	Role       string  `json:"role,omitempty"`
	FunctionId string  `json:"functionId,omitempty"`
	RequestId  string  `json:"request_id,omitempty"`
}

func NewMsgRollCallResponse(code string, role string) *MsgRollCallResponse {
	return &MsgRollCallResponse{
		Type: enums.MsgRollCallResponse,
		Code: code,
		Role: role,
	}
}

type MsgInstallFunction struct {
	Type        string  `json:"type,omitempty"`
	From        peer.ID `json:"from,omitempty"`
	ManifestUrl string  `json:"manifestUrl,omitempty"`
	Cid         string  `json:"cid,omitempty"`
}

func NewMsgInstallFunction(manifestUrl string) *MsgInstallFunction {
	return &MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: manifestUrl,
	}
}

type MsgInstallFunctionResponse struct {
	Type    string  `json:"type,omitempty"`
	From    peer.ID `json:"from,omitempty"`
	Code    string  `json:"code,omitempty"`
	Message string  `json:"message,omitempty"`
}

func NewMsgInstallFunctionResponse(code string, message string) *MsgInstallFunctionResponse {
	return &MsgInstallFunctionResponse{
		Type:    enums.MsgInstallFunctionResponse,
		Code:    code,
		Message: message,
	}
}
