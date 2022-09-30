package controller

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/executor"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func ExecuteFunction(ctx context.Context) string {
	out, err := executor.Execute(ctx)

	if err != nil {
		return "failed to execute function"
	}

	return string(out)
}

func RollCall(ctx context.Context) {
	msgRollCall := models.NewMsgRollCall()
	messaging.SendMessage(ctx, ctx.Value("topic").(*pubsub.Topic), msgRollCall)
}

func HealthStatus(ctx context.Context) {
	message := models.NewMsgHealthPing(enums.ResponseCodeOk)
	messaging.SendMessage(ctx, ctx.Value("topic").(*pubsub.Topic), message)
}
