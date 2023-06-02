package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	ma "github.com/multiformats/go-multiaddr"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	listenF := pflag.IntP("listen", "l", 0, "wait for incoming connections")
	insecureF := pflag.BoolP("insecure", "i", false, "use an unencrypted connection")
	privKeyPathF := pflag.StringP("private-key", "p", "", "path to the private key file")
	allowedPeerF := pflag.StringP("allowed-peer", "a", "", "allowed peer ID")
	pflag.Parse()

	if *listenF == 0 {
		logger.Fatal().Msg("Please provide a port to bind on with -l")
	}

	if *privKeyPathF == "" {
		logger.Fatal().Msg("Please provide a path to the private key file with -p")
	}

	ha, err := makeBasicHost(*listenF, *insecureF, *privKeyPathF)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create host")
	}

	startListener(ctx, ha, *listenF, *insecureF, *allowedPeerF, logger)
	<-ctx.Done()
}


func makeBasicHost(listenPort int, insecure bool, privKeyPath string) (host.Host, error) {
	priv, err := loadPrivateKey(privKeyPath)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	if insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	return libp2p.New(opts...)
}



func loadPrivateKey(filePath string) (crypto.PrivKey, error) {
    keyBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    priv, err := crypto.UnmarshalPrivateKey(keyBytes)
    if err != nil {
        return nil, err
    }

    return priv, nil
}

func getHostAddress(ha host.Host) string {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().Pretty()))
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func startListener(ctx context.Context, ha host.Host, listenPort int, insecure bool, allowedPeer string, logger zerolog.Logger) {
	fullAddr := getHostAddress(ha)
	logger.Info().Msgf("I am %s", fullAddr)

	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		if allowedPeer == "" && s.Conn().RemotePeer().Pretty() != allowedPeer {
			logger.Info().Msg("Connection from disallowed peer")
			s.Reset()
			return
		}
		logger.Info().Msg("listener received new stream")
		if err := handleRequest(s); err != nil {
			logger.Info().Err(err).Msg("Error in doEcho")
			s.Reset()
		} else {
			s.Close()
		}
	})

	logger.Info().Msg("listening for connections")
}

func handleRequest(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	var msg Message
	err = json.Unmarshal([]byte(str), &msg)
	if err != nil {
		return err
	}

	baseURL := "https://github.com/blocklessnetwork/b7s/releases/download"
	version := "v0.0.25"

	runtimeBaseURL := "https://github.com/blocklessnetwork/runtime/releases/download"
	runtimeVersion := "v0.0.12"

	responseMsg := &Message{}

	switch msg.Type {
	case "install_bls":
		go func() {
			installB7s(baseURL, version)
			installRuntime(runtimeBaseURL, runtimeVersion)
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			binPath := filepath.Join(usr.HomeDir, ".b7s", "bin")
			createServiceAndStartB7s(binPath)
		}()

	case "queryProcess":
		processInfo, err := CheckB7sRunning("b7s")
		if err != nil {
			log.Printf("Error checking b7s process: %v", err)
		} else if processInfo != nil {
			log.Printf("b7s process is running (PID: %d, Cmdline: %s)", processInfo.Pid, processInfo.Cmdline)
		} else {
			log.Println("b7s process is not running")
		}
	}

	responseData, _ := json.Marshal(responseMsg)
	_, err = s.Write(append(responseData, '\n'))
	return err
}
