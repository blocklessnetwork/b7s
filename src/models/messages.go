package models

import (
	"github.com/blocklessnetworking/b7s/src/enums"
)

type MsgBase struct {
	Type string `json:"type,omitempty"`
}

type MsgHealthPing struct {
	Type string `json:"type,omitempty"`
	Code string `json:"code,omitempty"`
}

func NewMsgHealthPing(code string) *MsgHealthPing {
	return &MsgHealthPing{
		Type: enums.MsgHealthCheck,
		Code: code,
	}
}

type MsgExecute struct {
	Type string `json:"type,omitempty"`
	Code string `json:"code,omitempty"`
}

func NewMsgExecute(code string) *MsgExecute {
	return &MsgExecute{
		Type: enums.MsgExecute,
		Code: code,
	}
}

type MsgExecuteResponse struct {
	Type   string `json:"type,omitempty"`
	Code   string `json:"code,omitempty"`
	Result string `json:"result,omitempty"`
}

type MsgRollCall struct {
	Type       string `json:"type,omitempty"`
	FunctionId string `json:"functionId,omitempty"`
}

func NewMsgRollCall(functionId string) *MsgRollCall {
	return &MsgRollCall{
		Type:       enums.MsgRollCall,
		FunctionId: functionId,
	}
}

type MsgRollCallResponse struct {
	Type string `json:"type,omitempty"`
	Code string `json:"code,omitempty"`
	Role string `json:"role,omitempty"`
}

func NewMsgRollCallResponse(code string, role string) *MsgRollCallResponse {
	return &MsgRollCallResponse{
		Type: enums.MsgRollCallResponse,
		Code: code,
		Role: role,
	}
}

type MsgInstallFunction struct {
	Type        string `json:"type,omitempty"`
	ManifestUrl string `json:"manifestUrl,omitempty"`
}

func NewMsgInstallFunction(manifestUrl string) *MsgInstallFunction {
	return &MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: manifestUrl,
	}
}
