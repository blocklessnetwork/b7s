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
)

func ExecuteFunction(ctx context.Context) (models.ExecutorResponse, error) {
	out, err := executor.Execute(ctx)

	if err != nil {
		return out, err
	}

	return out, nil
}

func InstallFunction(ctx context.Context, manifestPath string) error {
	repository.GetPackage(ctx, manifestPath)
	return nil
}

func RollCall(ctx context.Context) {
	msgRollCall := models.NewMsgRollCall()
	messaging.SendMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msgRollCall)
}

func HealthStatus(ctx context.Context) {
	message := models.NewMsgHealthPing(enums.ResponseCodeOk)
	messaging.SendMessage(ctx, ctx.Value("topic").(*pubsub.Topic), message)
}

func GetExecutionResponse(ctx context.Context, reqId string) *models.MsgExecuteResponse {
	memstore := ctx.Value("executionResponseMemStore").(memstore.ReqRespStore)
	response := memstore.Get(reqId)
	return response
}
