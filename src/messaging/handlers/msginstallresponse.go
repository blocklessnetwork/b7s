package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func HandleMsgInstallResponse(ctx context.Context, message []byte) {
	//log hello
	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message to install")
}
