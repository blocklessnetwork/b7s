package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/blocklessnetwork/orchestration-chain/x/market/types"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/libp2p/go-libp2p-core/host"
	log "github.com/sirupsen/logrus"
)

func registerNode(ctx context.Context) {
	var cfg = ctx.Value("config").(models.Config)

	log.Info("starting chain client service")

	var chainClientPath string
	addressPrefix := "bls"

	host := ctx.Value("host").(host.Host)

	userHomeDir, err := os.UserHomeDir()
	accountName := cfg.Chain.AddressKey

	if err != nil {
		log.Warn(err)
		return
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
		log.Warn(err)
		return
	}

	account, err := cosmosclient.Account(accountName)
	if err != nil {
		log.Warn(err)
		return
	}
	address, _ := account.Address(addressPrefix)
	msg := &types.MsgRegisterHeadNode{
		Creator:   address,
		NodeId:    host.ID().Pretty(),
		NodePort:  cfg.Node.Port,
		NodeIp:    cfg.Node.IP,
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
			log.Warn(err)
		}

		log.WithFields(log.Fields{
			"tx":        txResp.TxHash,
			"NodeOwner": address,
			"NodeId":    host.ID().Pretty(),
		}).Info("[chainservice] node registered")
	}
}

func startClient(ctx context.Context) {

}

func Start(ctx context.Context) {
	registerNode(ctx)
	go startClient(ctx)
}
