package handlers

import (
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

// roll call is recieved from pubsub peer
// respond to the requestor of the roll call directly
func HandleMsgRollCall(ctx context.Context, message []byte) {
	cfg := ctx.Value("config").(models.Config)

	// right now only workers should respond to roll calls
	// todo : other nodes should be able to respond to roll calls
	if cfg.Protocol.Role != enums.RoleWorker {
		return
	}

	msgRollCall := &models.MsgRollCall{}
	json.Unmarshal(message, msgRollCall)
	msgRollCall.From = ctx.Value("peerID").(peer.ID)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("rollcall message")

	channel := ctx.Value(enums.ChannelMsgLocal).(chan models.Message)

	localMsg := models.Message{
		Type: enums.MsgRollCall,
		Data: *msgRollCall,
	}

	channel <- localMsg
}

func HandleMsgRollCallResponse(ctx context.Context, message []byte) {
	msgRollCallResponse := &models.MsgRollCallResponse{}
	json.Unmarshal(message, msgRollCallResponse)
	msgRollCallResponse.From = ctx.Value("peerID").(peer.ID)

	log.WithFields(log.Fields{
		"message": string(message),
	}).Info("rollcall response")

	rollcallResponseChannel := ctx.Value(enums.ChannelMsgRollCallResponse).(chan models.MsgRollCallResponse)
	rollcallResponseChannel <- *msgRollCallResponse
}
