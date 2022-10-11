package handlers

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func HandleMsgHealthCheck(ctx context.Context, message []byte) {
	log.WithFields(log.Fields{
		"message": string(message),
	}).Debug("peer health check recieved")
}
