package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func HandleMsgExecute(ctx context.Context, message []byte) {
	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message from peer")
}
