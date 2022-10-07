package messaging

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

func Subscribe(ctx context.Context, host host.Host, topicName string) *pubsub.Topic {
	// make sure we're subscribed to the topic before we start publishing
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	topic, err := ps.Join(topicName)
	if err != nil {
		panic(err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, "topic", topic)
	// listen to messages
	go ListenMessages(ctx, sub, host)
	return topic
}

func SendMessage(ctx context.Context, topic *pubsub.Topic, message any) {
	messageString, _ := json.Marshal(message)
	if err := topic.Publish(ctx, []byte(messageString)); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("message err")
	}
}

func ListenMessages(ctx context.Context, sub *pubsub.Subscription, host host.Host) {
	for {
		message, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		if message.ReceivedFrom != host.ID() {
			HandleMessage(ctx, message)
		}
	}
}
