package controller

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/executor"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/repository"
	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func IsFunctionInstalled(ctx context.Context, functionId string) (models.FunctionManifest, error) {
	appDb := ctx.Value("appDb").(*pebble.DB)
	functionManifestString, err := db.Value(appDb, functionId)
	functionManifest := models.FunctionManifest{}

	json.Unmarshal([]byte(functionManifestString), &functionManifest)

	if err != nil {
		if err.Error() == "pebble: not found" {
			return functionManifest, errors.New("function not installed")
		} else {
			return functionManifest, err
		}
	}

	return functionManifest, nil
}

func ExecuteFunction(ctx context.Context, request models.RequestExecute) (models.ExecutorResponse, error) {
	config := ctx.Value("config").(models.Config)
	executorRole := config.Protocol.Role
	var out models.ExecutorResponse

	// if the role is a peer, then we need to send the request to the peer
	if executorRole == enums.RoleWorker {
		functionManifest, err := IsFunctionInstalled(ctx, request.FunctionId)

		// return if the function isn't installed
		// maybe install it?
		if err != nil {
			out := models.ExecutorResponse{
				Code: enums.ResponseCodeNotFound,
			}
			return out, err
		}

		out, err := executor.Execute(ctx, request, functionManifest)

		if err != nil {
			return out, err
		}

		return out, nil
	} else {
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
			out := models.ExecutorResponse{
				Code: enums.ResponseCodeTimeout,
			}
			return out, nil
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
			return out, err
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
}

func MsgInstallFunction(ctx context.Context, manifestPath string) {
	msg := models.MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: manifestPath,
	}

	log.Debug("request to message peer for install function")
	messaging.PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msg)
}

func InstallFunction(ctx context.Context, manifestPath string) error {
	repository.GetPackage(ctx, manifestPath)
	return nil
}

func RollCall(ctx context.Context, functionId string) *models.MsgRollCall {
	msgRollCall := models.NewMsgRollCall(functionId)
	messaging.PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msgRollCall)
	return msgRollCall
}

func RollCallResponse(ctx context.Context, msg models.MsgRollCall) {
	_, err := IsFunctionInstalled(ctx, msg.FunctionId)

	response := models.MsgRollCallResponse{
		Type:       enums.MsgRollCallResponse,
		FunctionId: msg.FunctionId,
		Code:       enums.ResponseCodeAccepted,
		RequestId:  msg.RequestId,
	}

	if err != nil {
		response.Code = enums.ResponseCodeNotFound
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		log.Debug("error marshalling roll call response")
		return
	}

	messaging.SendMessage(ctx, msg.From, responseJSON)
}

func HealthStatus(ctx context.Context) {
	message := models.NewMsgHealthPing(enums.ResponseCodeOk)
	messaging.PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), message)
}

func GetExecutionResponse(ctx context.Context, reqId string) *models.MsgExecuteResponse {
	memstore := ctx.Value("executionResponseMemStore").(memstore.ReqRespStore)
	response := memstore.Get(reqId)
	return response
}
