package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/messaging/handlers"
	"github.com/blocklessnetworking/b7s/src/models"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func HandleMessage(ctx context.Context, message []byte, peerID peer.ID) {
	var msg models.MsgBase
	ctx = context.WithValue(ctx, "peerID", peerID)
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Errorf("error unmarshalling message: %v", err)
	}

	var response interface{}

	handlers := map[string]func(context.Context, []byte){
		enums.MsgHealthCheck:      handlers.HandleMsgHealthCheck,
		enums.MsgExecute:          handlers.HandleMsgExecute,
		enums.MsgExecuteResponse:  handlers.HandleMsgExecuteResponse,
		enums.MsgRollCall:         handlers.HandleMsgRollCall,
		enums.MsgRollCallResponse: handlers.HandleMsgRollCallResponse,
		enums.MsgInstallFunction:  handlers.HandleMsgInstall,
	}

	if handler, ok := handlers[msg.Type]; ok {
		handler(ctx, message)
	}

	if topic := ctx.Value("topic").(*pubsub.Topic); topic != nil {
		PublishMessage(ctx, topic, response)
	}
}
