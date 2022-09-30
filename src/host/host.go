package host

import (
	"context"
	"strconv"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
)

func NewHost(ctx context.Context, port int, address string) host.Host {
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/" + address + "/tcp/" + strconv.FormatInt(int64(port), 10)))
	if err != nil {
		panic(err)
	}
	return host
}
