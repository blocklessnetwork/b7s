package node

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/response"
	"github.com/blocklessnetworking/b7s/node/internal/cache"
)

// TODO: Consider - interface for libp2p host instead of a type?

// TODO: Add doc comment.
type Node struct {
	role      blockless.NodeRole
	topicName string

	log      zerolog.Logger
	host     *host.Host
	store    Store
	execute  Execute
	function Function
	excache  *cache.Cache
	handlers map[string]HandlerFunc

	topic *pubsub.Topic

	rollCallResponses map[string](chan response.RollCall)
	executeResponses  map[string](chan response.Execute)
}

// New creates a new Node.
func New(log zerolog.Logger, host *host.Host, store Store, execute Execute, peerStore PeerStore, function Function, options ...func(*Config)) (*Node, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	n := Node{
		role:      cfg.Role,
		topicName: cfg.Topic,
		excache:   cache.New(),

		log:      log,
		host:     host,
		store:    store,
		function: function,
		execute:  execute,

		rollCallResponses: make(map[string](chan response.RollCall)),
		executeResponses:  make(map[string](chan response.Execute)),
	}

	// Initialize a list of handlers.
	handlers := map[string]HandlerFunc{
		blockless.MessageHealthCheck:             n.processHealthCheck,
		blockless.MessageExecute:                 n.getProcessHandlerFunc(n.role),
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
