package handlers

import (
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

func HandleMsgInstall(ctx context.Context, message []byte) {
	msgInstall := &models.MsgInstallFunction{}
	json.Unmarshal(message, msgInstall)
	msgInstall.From = ctx.Value("peerID").(peer.ID)

	channel := ctx.Value(enums.ChannelMsgInstallFunction).(chan models.MsgInstallFunction)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message to install")

	channel <- *msgInstall
}
