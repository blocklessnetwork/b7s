package messaging

import (
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/enums"
	handlers "github.com/blocklessnetworking/b7s/src/messaging/handlers"
	"github.com/blocklessnetworking/b7s/src/models"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func HandleMessage(ctx context.Context, message []byte, peerID peer.ID) {
	var msg models.MsgBase
	ctx = context.WithValue(ctx, "peerID", peerID)
	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		panic(err)
	}

	var response interface{}
	switch msg.Type {
	case enums.MsgHealthCheck:
		handlers.HandleMsgHealthCheck(ctx, message)
	case enums.MsgExecute:
		handlers.HandleMsgExecute(ctx, message)
	case enums.MsgRollCall:
		handlers.HandleMsgRollCall(ctx, message)
	case enums.MsgRollCallResponse:
		handlers.HandleMsgRollCallResponse(ctx, message)
	case enums.MsgInstallFunction:
		handlers.HandleMsgInstall(ctx, message)
	}

	if response != nil {
		PublishMessage(ctx, ctx.Value("topic").(*pubsub.Topic), response)
	}
}
