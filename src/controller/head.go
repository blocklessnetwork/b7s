package controller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
)

func HeadExecuteFunction(ctx context.Context, request models.RequestExecute) (models.ExecutorResponse, error) {
	// perform rollcall to see who is available
	rollcallMessage := RollCall(ctx, request.FunctionId)

	type rollcalled struct {
		From peer.ID
		Code string
	}
	rollCalledChannel := make(chan rollcalled)

	go func(ctx context.Context) {
		// collect responses of nodes who want to work on the request
		host := ctx.Value("host").(host.Host)
		rollcallResponseChannel := ctx.Value(enums.ChannelMsgRollCallResponse).(chan models.MsgRollCallResponse)
		_, timeoutCancel := context.WithCancel(ctx)
		// time out
		// should we retry this?
		// possible no ne responds back
		go func() {
			timeout := time.After(5 * time.Second)
		LOOP:
			select {
			case <-timeout:
				rollCalledChannel <- rollcalled{
					Code: enums.ResponseCodeTimeout,
				}
				return
			case msg := <-rollcallResponseChannel:
				conns := host.Network().ConnsToPeer(msg.From)

				// pop off all the responses that don't match our first found connection
				// worker with function
				// worker responsed accepted has resources and is ready to execute
				// worker knows RequestID
				if msg.Code == enums.ResponseCodeAccepted && msg.FunctionId == request.FunctionId && len(conns) > 0 && rollcallMessage.RequestId == msg.RequestId {
					log.WithFields(log.Fields{
						"msg": msg,
					}).Info("rollcalled")
					timeoutCancel()
					rollCalledChannel <- rollcalled{
						From: msg.From,
					}
					return
				} else {
					goto LOOP
				}
			case <-ctx.Done():
				log.Warn("timeout cancelled")
				return
			}
		}()
	}(ctx)

	msgRollCall := <-rollCalledChannel

	if msgRollCall.Code == enums.ResponseCodeTimeout {
		return models.ExecutorResponse{
			Code: enums.ResponseCodeTimeout,
		}, nil
	}

	// we got a message back before the timeout went off
	// request an execution from first responding node
	// we should queue these responses into a pool first
	// for selection
	msgExecute := models.MsgExecute{
		Type:       enums.MsgExecute,
		FunctionId: request.FunctionId,
		Method:     request.Method,
		Parameters: request.Parameters,
		Config:     request.Config,
	}

	jsonBytes, err := json.Marshal(msgExecute)

	if err != nil {
		return models.ExecutorResponse{
			Code: enums.ResponseCodeError,
		}, err
	}

	// send execute message to node
	messaging.SendMessage(ctx, msgRollCall.From, jsonBytes)

	// wait for response
	executeResponseChannel := ctx.Value(enums.ChannelMsgExecuteResponse).(chan models.MsgExecuteResponse)
	msgExec := <-executeResponseChannel

	// too many models here ?
	out := models.ExecutorResponse{
		Code:      msgExec.Code,
		Result:    msgExec.Result,
		RequestId: msgExec.RequestId,
	}

	log.WithFields(log.Fields{
		"msg": msgExec,
	}).Info("execute response")

	defer close(rollCalledChannel)
	return out, nil

}
