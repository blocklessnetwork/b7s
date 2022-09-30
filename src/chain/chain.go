package chain

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blocklessnetwork/orchestration-chain/x/market/types"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/libp2p/go-libp2p-core/host"
	log "github.com/sirupsen/logrus"
)

func registerNode(ctx context.Context) {
	var chainClientPath string
	addressPrefix := "bls"
	cfg, _ := ctx.Value("config").(*models.Config)
	host := ctx.Value("host").(host.Host)

	userHomeDir, err := os.UserHomeDir()
	accountName := cfg.Chain.AddressKey

	if err != nil {
		panic(err)
	}

	if len(cfg.Chain.Home) > 0 {
		chainClientPath = cfg.Chain.Home
	} else {
		chainClientPath = filepath.Join(userHomeDir, ".blockless-chain")
	}

	chainOptions := []cosmosclient.Option{
		cosmosclient.WithKeyringBackend("test"),
		cosmosclient.WithHome(chainClientPath),
		cosmosclient.WithAddressPrefix(addressPrefix),
		cosmosclient.WithNodeAddress(cfg.Chain.RPC),
	}

	cosmosclient, err := cosmosclient.New(ctx, chainOptions...)

	if err != nil {
		log.Fatal(err)
	}

	account, err := cosmosclient.Account(accountName)
	if err != nil {
		log.Fatal(err)
	}
	address, _ := account.Address(addressPrefix)
	msg := &types.MsgRegisterHeadNode{
		Creator:   address,
		NodeId:    host.ID().Pretty(),
		NodePort:  strconv.Itoa(cfg.Node.Port),
		NodeIp:    cfg.Node.IpAddress,
		NodeOwner: address,
	}
	queryClient := types.NewQueryClient(cosmosclient.Context())
	registrations, _ := queryClient.NodeRegistration(context.Background(), &types.QueryGetNodeRegistrationRequest{
		Index: address + "-" + host.ID().Pretty(),
	})

	if registrations != nil && len(registrations.NodeRegistration.NodeId) > 0 {
		log.WithFields(log.Fields{
			"NodeOwner": address,
			"NodeId":    address + "-" + host.ID().Pretty(),
		}).Info("[chainservice] node already registered ")
	} else {
		txResp, err := cosmosclient.BroadcastTx(account, msg)
		if err != nil {
			log.Fatal(err)
		}

		log.WithFields(log.Fields{
			"tx":        txResp.TxHash,
			"NodeOwner": address,
			"NodeId":    host.ID().Pretty(),
		}).Info("[chainservice] node registered")
	}
}

func Start(ctx context.Context) {
	// var config = ctx.Value("config").(models.Config)

	log.Info("starting blockchain client")
	go registerNode(ctx)
}
