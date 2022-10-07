package host

import (
	"context"
	"strconv"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
)

func NewHost(ctx context.Context, port int, address string) host.Host {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/" + address + "/tcp/" + strconv.FormatInt(int64(port), 10)),
		// libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}
	return host
}
