package host

import (
	"context"
	"io/ioutil"
	"strconv"

	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	log "github.com/sirupsen/logrus"
)

// NewHost creates a new libp2p host
func NewHost(ctx context.Context, port int, address string) host.Host {
	var privKey crypto.PrivKey

	// Read the private key file if it exists
	keyPath := ctx.Value("config").(models.Config).Node.KeyPath
	if keyPath != "" {
		log.Println("loading private key from: ", keyPath)
		privKeyBytes, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Error("failed to load private key from: ", keyPath)
		}

		key, err := crypto.UnmarshalPrivateKey(privKeyBytes)
		if err != nil {
			log.Error("failed to load private key from: ", keyPath)
		}
		privKey = key
	}

	var hostAddress = "/ip4/" + address + "/tcp/" + strconv.FormatInt(int64(port), 10)
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(hostAddress),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// Use the private key if it exists, otherwise generate an identity when starting the host
	if privKey != nil {
		opts = append(opts, libp2p.Identity(privKey))
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}

	log.Info("host: ", hostAddress+"/p2p/"+h.ID().Pretty())

	// Set a stream handler to listen for incoming streams
	messaging.ListenMessages(ctx, h)

	return h
}
