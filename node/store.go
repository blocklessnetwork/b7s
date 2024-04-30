package node

import (
	"github.com/blocklessnetwork/b7s/models/blockless"
)

type Store interface {
	SavePeer(blockless.Peer) error
}
