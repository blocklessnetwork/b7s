package messaging

import (
	"bufio"
	"context"
	"encoding/json"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	log "github.com/sirupsen/logrus"
)

// subscribe to a gossipsub topic
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
	go ListenPublishedMessages(ctx, sub, host)
	return topic
}

// publish messages on the gossipsub topic
func PublishMessage(ctx context.Context, topic *pubsub.Topic, message any) {
	messageString, _ := json.Marshal(message)
	if err := topic.Publish(ctx, []byte(messageString)); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("message err")
	}
}

// listens to pubsub messages and send them to a message handler
func ListenPublishedMessages(ctx context.Context, sub *pubsub.Subscription, host host.Host) {
	for {
		message, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		if message.ReceivedFrom != host.ID() {
			HandleMessage(ctx, message.Data, message.ReceivedFrom)
		}
	}
}

// listen to direct messages from peers
func ListenMessages(ctx context.Context, host host.Host) {
	host.SetStreamHandler(enums.WorkProtocolId, func(s network.Stream) {
		defer s.Close()
		buf := bufio.NewReader(s)
		response, err := buf.ReadString('\n')

		if err != nil {
			s.Reset()
			log.Warn(err)
		}

		HandleMessage(ctx, []byte(response), s.Conn().RemotePeer())
	})
}

// sends a message directly to a peer
func SendMessage(ctx context.Context, peer peer.ID, message []byte) {
	host := ctx.Value("host").(host.Host)
	s, err := host.NewStream(context.Background(), peer, enums.WorkProtocolId)
	if err != nil {
		log.Warn(err)
	}
	_, err = s.Write(message)
	if err != nil {
		s.Reset()
		log.Warn(err)
	}
	defer s.Close()
}
