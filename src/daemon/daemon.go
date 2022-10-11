package daemon

import (
	"context"
	"strconv"

	"os"
	"path/filepath"
	"time"

	"github.com/blocklessnetworking/b7s/src/chain"
	"github.com/blocklessnetworking/b7s/src/config"
	"github.com/blocklessnetworking/b7s/src/controller"
	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/dht"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/health"
	"github.com/blocklessnetworking/b7s/src/host"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/restapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// the daemonm service loop
// also the rootCommand for cobra
func Run(cmd *cobra.Command, args []string, configPath string) {
	topicName := "blockless.networking/networking/general"
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

	// define channels before instanciating the host
	msgInstallFunctionChannel := make(chan models.MsgInstallFunction)
	msgRollCallChannel := make(chan models.MsgRollCall)
	msgRollCallResponseChannel := make(chan models.MsgRollCallResponse)
	ctx = context.WithValue(ctx, enums.ChannelMsgInstallFunction, msgInstallFunctionChannel)
	ctx = context.WithValue(ctx, enums.ChannelMsgRollCall, msgRollCallChannel)
	ctx = context.WithValue(ctx, enums.ChannelMsgRollCallResponse, msgRollCallResponseChannel)

	host := host.NewHost(ctx, port, config.C.Node.IP)
	ctx = context.WithValue(ctx, "host", host)

	// set appdb config
	appDb := db.Get(exPath + "/" + host.ID().Pretty() + "_appDb")
	ctx = context.WithValue(ctx, "appDb", appDb)

	// response memstore
	// todo flush memstore occasionally
	executionResponseMemStore := memstore.NewReqRespStore()
	ctx = context.WithValue(ctx, "executionResponseMemStore", executionResponseMemStore)

	go (func() {
		for {
			select {
			case msg := <-msgInstallFunctionChannel:
				controller.InstallFunction(ctx, msg.ManifestUrl)
			case msg := <-msgRollCallChannel:
				controller.RollCallResponse(ctx, msg)
			}
		}
	})()

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

	// run the daemon loop
	select {}
}
