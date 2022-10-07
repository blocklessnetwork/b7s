package handlers

import (
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

func HandleMsgInstall(ctx context.Context, message []byte) {
	msgInstall := &models.MsgInstallFunction{}
	json.Unmarshal(message, msgInstall)
	channel := ctx.Value("channelHandler").(chan models.MsgInstallFunction)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message from peer")

	channel <- *msgInstall
}
