package host

import (
	"bufio"
	"context"
	"io/ioutil"

	"strconv"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	log "github.com/sirupsen/logrus"
)

func NewHost(ctx context.Context, port int, address string) host.Host {

	// see if we have a private key to load
	var privKey crypto.PrivKey
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

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/" + address + "/tcp/" + strconv.FormatInt(int64(port), 10)),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// load private key if we have one
	// otherwise we will just generate an identity when we start the host
	if privKey != nil {
		opts = append(opts, libp2p.Identity(privKey))
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}

	// set a stream handler on the worker to listen for incoming streams
	// from a head node
	if ctx.Value("config").(models.Config).Protocol.Role == "worker" {
		go func() {
			host.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
				log.Println("listener received new stream")
				if err := doEcho(s); err != nil {
					log.Println(err)
					s.Reset()
				} else {
					s.Close()
				}
			})
		}()
	}

	return host
}

func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
