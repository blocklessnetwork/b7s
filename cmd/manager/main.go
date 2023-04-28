package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/spf13/pflag"

	golog "github.com/ipfs/go-log/v2"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	golog.SetAllLoggers(golog.LevelInfo)

	listenF := pflag.IntP("listen", "l", 0, "wait for incoming connections")
	insecureF := pflag.BoolP("insecure", "i", false, "use an unencrypted connection")
	seedF := pflag.Int64P("seed", "s", 0, "set random seed for id generation")
	pflag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	ha, err := makeBasicHost(*listenF, *insecureF, *seedF)
	if err != nil {
		log.Fatal(err)
	}

	startListener(ctx, ha, *listenF, *insecureF)
	<-ctx.Done()
}

func makeBasicHost(listenPort int, insecure bool, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
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

func getHostAddress(ha host.Host) string {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().Pretty()))
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func startListener(ctx context.Context, ha host.Host, listenPort int, insecure bool) {
	fullAddr := getHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	log.Println("listening for connections")

	if insecure {
		log.Printf("Now run \"./echo -l %d -d %s --insecure\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"./echo -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}
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