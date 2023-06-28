package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/node/internal/waitmap"
)

// Node is the entity that actually provides the main Blockless node functionality.
// It listens for messages coming from the wire and processes them. Depending on the
// node role, which is determined on construction, it may process messages in different ways.
// For example, upon receiving a message requesting execution of a Blockless function,
// a Worker Node will use the `Execute` component to fullfill the execution request.
// On the other hand, a Head Node will issue a roll call and eventually
// delegate the execution to the chosend Worker Node.
type Node struct {
	cfg Config

	log      zerolog.Logger
	host     *host.Host
	executor Executor
	fstore   FStore

	topic *pubsub.Topic
	sema  chan struct{}
	wg    *sync.WaitGroup

	rollCall *rollCallQueue

	// clusters maps request ID to the raft cluster the node belongs to.
	// TODO: (raft) think when the raft cluster should be disbanded and this removed.
	clusters map[string]*raftHandler

	// clusterLock is used to synchronize access to the `clusters` map.
	clusterLock sync.RWMutex

	executeResponses   *waitmap.WaitMap
	consensusResponses *waitmap.WaitMap
}

// New creates a new Node.
func New(log zerolog.Logger, host *host.Host, peerStore PeerStore, fstore FStore, options ...Option) (*Node, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	n := &Node{
		cfg: cfg,

		log:      log.With().Str("component", "node").Logger(),
		host:     host,
		fstore:   fstore,
		executor: cfg.Execute,

		wg:   &sync.WaitGroup{},
		sema: make(chan struct{}, cfg.Concurrency),

		rollCall:           newQueue(rollCallQueueBufferSize),
		clusters:           make(map[string]*raftHandler),
		executeResponses:   waitmap.New(),
		consensusResponses: waitmap.New(),
	}

	err := n.ValidateConfig()
	if err != nil {
		return nil, fmt.Errorf("node configuration is not valid: %w", err)
	}

	// Create a notifiee with a backing peerstore.
	cn := newConnectionNotifee(log, peerStore)
	host.Network().Notify(cn)

	return n, nil
}

// ID returns the ID of this node.
func (n *Node) ID() string {
	return n.host.ID().String()
}

// getHandler returns the appropriate handler function for the given message.
func (n *Node) getHandler(msgType string) HandlerFunc {

	switch msgType {
	case blockless.MessageHealthCheck:
		return n.processHealthCheck
	case blockless.MessageExecuteResponse:
		return n.processExecuteResponse
	case blockless.MessageRollCall:
		return n.processRollCall
	case blockless.MessageRollCallResponse:
		return n.processRollCallResponse
	case blockless.MessageInstallFunction:
		return n.processInstallFunction
	case blockless.MessageInstallFunctionResponse:
		return n.processInstallFunctionResponse
	case blockless.MessageFormCluster:
		return n.processFormCluster
	case blockless.MessageFormClusterResponse:
		return n.processFormClusterResponse
	case blockless.MessageDisbandCluster:
		return n.processDisbandCluster

	case blockless.MessageExecute:
		return n.processExecute

	default:
		return func(_ context.Context, from peer.ID, _ []byte) error {
			return ErrUnsupportedMessage
		}
	}
}

func newRequestID() (string, error) {

	// Generate a new request/executionID.
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("could not generate new request ID: %w", err)
	}

	return uuid.String(), nil
}
