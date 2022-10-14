package daemon

import (
	"context"
	"strconv"

	"os"
	"path/filepath"
	"time"

	"github.com/blocklessnetworking/b7s/src/chain"
	"github.com/blocklessnetworking/b7s/src/config"
	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/dht"
	"github.com/blocklessnetworking/b7s/src/health"
	"github.com/blocklessnetworking/b7s/src/host"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/restapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// the daemonm service loop
// also the rootCommand for cobra
func Run(cmd *cobra.Command, args []string, configPath string) {
	topicName := "blockless/b7s/general"
	ctx := context.Background()
	ex, err := os.Executable()
	if err != nil {
		log.Warn(err)
	}

	// get the path to the executable
	exPath := filepath.Dir(ex)

	// load config
	err = config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// set context config
	ctx = context.WithValue(ctx, "config", config.C)

	// create a new node hode
	port, err := strconv.Atoi(config.C.Node.Port)
	if err != nil {
		log.Fatal(err)
	}

	// setup any daemon processing channels we need to listen to
	ctx = setupChannels(ctx)

	// create a new libp2p host
	host := host.NewHost(ctx, port, config.C.Node.IP)
	ctx = context.WithValue(ctx, "host", host)

	// set appdb config
	appDb := db.Get(exPath + "/" + host.ID().Pretty() + "_appDb")
	ctx = context.WithValue(ctx, "appDb", appDb)

	// response memstore
	// todo flush memstore occasionally
	executionResponseMemStore := memstore.NewReqRespStore()
	ctx = context.WithValue(ctx, "executionResponseMemStore", executionResponseMemStore)

	// listen to channels and handle messages from other parts of the applications
	// if coming from the network, these messages would processing through `messageHandlers` first
	go listenToChannels(ctx)

	// pubsub topics from p2p
	topic := messaging.Subscribe(ctx, host, topicName)
	ctx = context.WithValue(ctx, "topic", topic)

	// start health monitoring
	ticker := time.NewTicker(1 * time.Minute)
	go health.StartPing(ctx, ticker)

	// start other services based on config
	if config.C.Protocol.Role == "head" {
		restapi.Start(ctx)
		chain.Start(ctx)
	}

	defer ticker.Stop()

	// discover peers
	go dht.DiscoverPeers(ctx, host, topicName)

	// daemon is running
	// waiting for ctrl-c to exit
	select {}
}
