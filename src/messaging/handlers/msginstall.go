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
	cfg := ctx.Value("config").(models.Config)

	// right now only workers should respond to install calls, so that
	// they are not particpating in work just yet
	if cfg.Protocol.Role != enums.RoleWorker {
		return
	}

	msgInstall := &models.MsgInstallFunction{}
	json.Unmarshal(message, msgInstall)
	msgInstall.From = ctx.Value("peerID").(peer.ID)

	channel := ctx.Value(enums.ChannelMsgLocal).(chan models.Message)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message to install")

	localMsg := models.Message{
		Type: enums.MsgInstallFunction,
		Data: *msgInstall,
	}

	channel <- localMsg
}
