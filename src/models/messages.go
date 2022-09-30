package models

import (
	"github.com/blocklessnetworking/b7s/src/enums"
)

type MsgBase struct {
	Type string `json:"type"`
}

type MsgHealthPing struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

func NewMsgHealthPing(code string) *MsgHealthPing {
	return &MsgHealthPing{
		Type: enums.MsgHealthCheck,
		Code: code,
	}
}

type MsgExecute struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

func NewMsgExecute(code string) *MsgExecute {
	return &MsgExecute{
		Type: enums.MsgExecute,
		Code: code,
	}
}

type MsgRollCall struct {
	Type string `json:"type"`
}

func NewMsgRollCall() *MsgRollCall {
	return &MsgRollCall{
		Type: enums.MsgRollCall,
	}
}

type MsgRollCallResponse struct {
	Type string `json:"type"`
	Code string `json:"code"`
	Role string `json:"role"`
}

func NewMsgRollCallResponse(code string, role string) *MsgRollCallResponse {
	return &MsgRollCallResponse{
		Type: enums.MsgRollCallResponse,
		Code: code,
		Role: role,
	}
}
