package handlers

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

func HandleMsgRollCall(ctx context.Context, message []byte) *models.MsgRollCallResponse {
	response := models.NewMsgRollCallResponse(enums.ResponseCodeOk, ctx.Value("config").(models.Config).Protocol.Role)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message from peer")

	return response
}
