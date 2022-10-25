package controller

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/repository"
	"github.com/cockroachdb/pebble"
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

	// if the role is a peer, then we need to send the request to the peer
	if executorRole == enums.RoleWorker {
		return WorkerExecuteFunction(ctx, request)
	} else {
		return HeadExecuteFunction(ctx, request)
	}
}

func MsgInstallFunction(ctx context.Context, installRequest models.RequestFunctionInstall) {
	var manifestPath string
	// determine the manifest path
	if installRequest.Uri != "" {
		manifestPath = installRequest.Uri
	}

	if installRequest.Cid != "" {
		// todo
		// use other portals when available
		manifestPath = "https://" + installRequest.Cid + ".ipfs.w3s.link/manifest.json"
	}

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
