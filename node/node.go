package node

import (
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// TODO: Consider - interface for libp2p host instead of a type?

// TODO: Add doc comment.
type Node struct {
	role     blockless.NodeRole
	topic    string
	handlers map[string]HandlerFunc

	log     zerolog.Logger
	host    *host.Host
	store   Store
	execute Execute
}

// New creates a new Node.
func New(log zerolog.Logger, host *host.Host, store Store, execute Execute, peerStore PeerStore, options ...func(*Config)) (*Node, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	n := Node{
		role:  cfg.Role,
		topic: cfg.Topic,

		log:     log,
		host:    host,
		store:   store,
		execute: execute,
	}

	// Initialize a list of handlers.
	handlers := map[string]HandlerFunc{
		blockless.MessageHealthCheck:             n.processHealthCheck,
		blockless.MessageExecute:                 n.processExecute,
		blockless.MessageExecuteResponse:         n.processExecuteResponse,
		blockless.MessageRollCall:                n.processRollCall,
		blockless.MessageRollCallResponse:        n.processRollCallResponse,
		blockless.MessageInstallFunction:         n.processInstallFunction,
		blockless.MessageInstallFunctionResponse: n.processInstallFunctionResponse,
	}
	n.handlers = handlers

	// Create a notifiee with a backing peerstore.
	cn := newConnectionNotifee(log, peerStore)
	host.Network().Notify(cn)

	return &n, nil
}
