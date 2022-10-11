package controller

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/executor"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/repository"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func ExecuteFunction(ctx context.Context, request models.RequestExecute, functionManifest models.FunctionManifest) (models.ExecutorResponse, error) {
	config := ctx.Value("config").(models.Config)
	executorRole := config.Protocol.Role
	var out models.ExecutorResponse

	// if the role is a peer, then we need to send the request to the peer
	if executorRole == enums.RoleWorker {
		out, err := executor.Execute(ctx, request, functionManifest)

		if err != nil {
			return out, err
		}

		return out, nil
	} else {
		// perform rollcall to see who is available
		RollCall(ctx, request.FunctionId)
	}

	return out, nil
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

func RollCall(ctx context.Context, functionId string) {
	msgRollCall := models.NewMsgRollCall(functionId)
	messaging.PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msgRollCall)
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
