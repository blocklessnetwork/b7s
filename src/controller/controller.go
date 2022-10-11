package controller

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/executor"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/repository"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func ExecuteFunction(ctx context.Context, request models.RequestExecute, functionManifest models.FunctionManifest) (models.ExecutorResponse, error) {
	out, err := executor.Execute(ctx, request, functionManifest)

	if err != nil {
		return out, err
	}

	return out, nil
}

func MsgInstallFunction(ctx context.Context, manifestPath string) {
	msg := models.MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: manifestPath,
	}

	log.Info("request to message peer for install function")
	host := ctx.Value("host").(host.Host)
	connectedPeers := host.Peerstore().PeersWithAddrs()
	for _, peer := range connectedPeers {
		if peer.Pretty() != host.ID().Pretty() {
			log.Info("sending install function message to peer: ", peer)
			s, err := host.NewStream(context.Background(), peer, "/echo/1.0.0")
			if err != nil {
				log.Warn(err)
			}
			_, err = s.Write([]byte("Hello, world!\n"))
			if err != nil {
				log.Warn(err)
			}
		}
	}

	messaging.PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msg)
}

func InstallFunction(ctx context.Context, manifestPath string) error {
	repository.GetPackage(ctx, manifestPath)
	return nil
}

func RollCall(ctx context.Context) {
	msgRollCall := models.NewMsgRollCall()
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
