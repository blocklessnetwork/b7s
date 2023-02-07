package node

import (
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/host"
)

// TODO: Add doc comment.
type Node struct {
	log zerolog.Logger

	// TODO: Check - interface for this instead of a type?
	host *host.Host
}

// New creates a new Node.
func New(log zerolog.Logger, host *host.Host, peerStore PeerStore) (*Node, error) {

	node := Node{
		log:  log,
		host: host,
	}

	// Create a notifiee with a backing peerstore.
	cn := newConnectionNotifee(log, peerStore)
	host.Notify(cn)

	return &node, nil
}
