package daemon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blocklessnetworking/b7s/src/controller"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

func setupChannels(ctx context.Context) context.Context {
	// define channels before instanciating the host
	// msgInstallFunctionChannel := make(chan models.MsgInstallFunction)
	// msgExecute := make(chan models.MsgExecute)
	msgExecuteResponse := make(chan models.MsgExecuteResponse)
	// msgRollCallChannel := make(chan models.MsgRollCall)
	msgRollCallResponseChannel := make(chan models.MsgRollCallResponse)
	msgChannelLocal := make(chan models.Message)
	// ctx = context.WithValue(ctx, enums.ChannelMsgExecute, msgExecute)
	ctx = context.WithValue(ctx, enums.ChannelMsgExecuteResponse, msgExecuteResponse)
	// ctx = context.WithValue(ctx, enums.ChannelMsgInstallFunction, msgInstallFunctionChannel)
	// ctx = context.WithValue(ctx, enums.ChannelMsgRollCall, msgRollCallChannel)
	ctx = context.WithValue(ctx, enums.ChannelMsgRollCallResponse, msgRollCallResponseChannel)
	ctx = context.WithValue(ctx, enums.ChannelMsgLocal, msgChannelLocal)

	return ctx
}

func listenToChannels(ctx context.Context) {

	msgChannel := ctx.Value(enums.ChannelMsgLocal).(chan models.Message)

	for {
		select {
		case msg := <-msgChannel:
			switch msg.Type {
			case enums.MsgInstallFunction:
				m, ok := msg.Data.(models.MsgInstallFunction)
				if ok {
					controller.InstallFunction(ctx, m)
				} else {
					fmt.Println("The assertion failed.")
				}
			case "execute":
				msg := msg.Data.(models.MsgExecute)
				requestExecute := models.RequestExecute{
					FunctionId: msg.FunctionId,
					Method:     msg.Method,
					Parameters: msg.Parameters,
					Config:     msg.Config,
				}
				executorResponse, err := controller.ExecuteFunction(ctx, requestExecute)
				if err != nil {
					log.Error(err)
				}
				jsonBytes, err := json.Marshal(&models.MsgExecuteResponse{
					RequestId: executorResponse.RequestId,
					Type:      enums.MsgExecuteResponse,
					Code:      executorResponse.Code,
					Result:    executorResponse.Result,
				})
				messaging.SendMessage(ctx, msg.From, jsonBytes)
			case enums.MsgRollCall:
				controller.RollCallResponse(ctx, msg.Data.(models.MsgRollCall))
			}
		}
	}
}

// func listenToChannels(ctx context.Context) {
// 	msgInstallFunctionChannel := ctx.Value(enums.ChannelMsgInstallFunction).(chan models.MsgInstallFunction)
// 	msgExecute := ctx.Value(enums.ChannelMsgExecute).(chan models.MsgExecute)
// 	msgRollCallChannel := ctx.Value(enums.ChannelMsgRollCall).(chan models.MsgRollCall)

// 	for {
// 		select {
// 		case msg := <-msgInstallFunctionChannel:
// 			controller.InstallFunction(ctx, msg)
// 		case msg := <-msgRollCallChannel:
// 			controller.RollCallResponse(ctx, msg)
// 		case msg := <-msgExecute:
// 			// todo no sir I don't like this
// 			// I think this is duplicated in the controller
// 			requestExecute := models.RequestExecute{
// 				FunctionId: msg.FunctionId,
// 				Method:     msg.Method,
// 				Parameters: msg.Parameters,
// 				Config:     msg.Config,
// 			}
// 			executorResponse, err := controller.ExecuteFunction(ctx, requestExecute)
// 			if err != nil {
// 				log.Error(err)
// 			}

// 			jsonBytes, err := json.Marshal(&models.MsgExecuteResponse{
// 				RequestId: executorResponse.RequestId,
// 				Type:      enums.MsgExecuteResponse,
// 				Code:      executorResponse.Code,
// 				Result:    executorResponse.Result,
// 			})

// 			// send exect response back to head node
// 			messaging.SendMessage(ctx, msg.From, jsonBytes)
// 		}
// 	}
// }
