package handlers

import (
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

// roll call is recieved from pubsub peer
// respond to the requestor of the roll call directly
func HandleMsgRollCall(ctx context.Context, message []byte) {
	msgRollCall := &models.MsgRollCall{}
	json.Unmarshal(message, msgRollCall)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("message from peer")

	channel := ctx.Value("msgRollCallChannel").(chan models.MsgRollCall)
	channel <- *msgRollCall
}
