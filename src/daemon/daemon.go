package daemon

import (
	"context"
	"log"
	"time"

	"github.com/blocklessnetworking/b7s/src/chain"
	"github.com/blocklessnetworking/b7s/src/config"
	"github.com/blocklessnetworking/b7s/src/dht"
	"github.com/blocklessnetworking/b7s/src/health"
	"github.com/blocklessnetworking/b7s/src/host"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/rest"
	"github.com/spf13/cobra"
)

// the daemonm service loop
// also the rootCommand for cobra
func Run(cmd *cobra.Command, args []string, configPath string) {
	topicName := "blockless.networking/networking/general"
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// set context config
	ctx = context.WithValue(ctx, "config", config.C)

	// create a new node hode
	host := host.NewHost(ctx, config.C.Node.Port, config.C.Node.IpAddress)
	ctx = context.WithValue(ctx, "host", host)

	// subscribe to public topic
	topic := messaging.Subscribe(ctx, host, topicName)
	ctx = context.WithValue(ctx, "topic", topic)

	// start health monitoring
	ticker := time.NewTicker(1 * time.Minute)
	go health.StartPing(ctx, ticker)

	// start other services
	rest.Start(ctx)
	chain.Start(ctx)

	defer ticker.Stop()

	// discover peers
	go dht.DiscoverPeers(ctx, host, topicName)

	// run the daemon loop
	select {}
}
