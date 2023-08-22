package raft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	libp2praft "github.com/libp2p/go-libp2p-raft"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// TODO (raft): Handler logging.

type Handler struct {
	*raft.Raft
	logStore *boltdb.BoltStore
	stable   *boltdb.BoltStore

	cfg Config
	log zerolog.Logger

	rootDir string
	peers   []peer.ID
}

// New creates a new raft handler, bootstraps the cluster and waits until a first leader is elected. We do this because
// only after the election the cluster is really operational and ready to process requests.
func New(log zerolog.Logger, host *host.Host, workspace string, requestID string, executor Executor, peers []peer.ID, options ...Option) (*Handler, error) {

	// Step 1: Create a new raft handler.
	h, err := newHandler(log, host, workspace, requestID, executor, peers, options...)
	if err != nil {
		return nil, fmt.Errorf("could not create raft handler: %w", err)
	}

	// Step 2: Register an observer to monitor leadership changes. More precisely,
	// wait on the first leader election, so we know when the cluster is operational.

	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, func(obs *raft.Observation) bool {
		_, ok := obs.Data.(raft.LeaderObservation)
		return ok
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait on leadership observation.
		obs := <-obsCh
		leaderObs, ok := obs.Data.(raft.LeaderObservation)
		if !ok {
			h.log.Error().Type("type", obs.Data).Msg("invalid observation type received")
			return
		}

		// We don't need the observer anymore.
		h.DeregisterObserver(observer)

		h.log.Info().Str("request", requestID).Str("leader", string(leaderObs.LeaderID)).Msg("observed a leadership event - ready")
	}()

	h.RegisterObserver(observer)

	// Step 3: Bootstrap the cluster.
	err = h.bootstrapCluster()
	if err != nil {
		return nil, fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	wg.Wait()

	return h, nil
}

func newHandler(log zerolog.Logger, host *host.Host, workspace string, requestID string, executor Executor, peers []peer.ID, options ...Option) (*Handler, error) {

	if len(peers) == 0 {
		return nil, errors.New("empty peer list")
	}

	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	// Determine directory that should be used for consensus for this request.
	rootDir := consensusDir(workspace, requestID)
	err := os.MkdirAll(rootDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create consensus work directory: %w", err)
	}

	// Transport layer for raft communication.
	transport, err := libp2praft.NewLibp2pTransport(host, consensusTransportTimeout)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p transport: %w", err)
	}

	// Create log store.
	logDB := filepath.Join(rootDir, defaultLogStoreName)
	logStore, err := boltdb.NewBoltStore(logDB)
	if err != nil {
		return nil, fmt.Errorf("could not create log store (path: %s): %w", logDB, err)
	}

	// Create stable store.
	stableDB := filepath.Join(rootDir, defaultStableStoreName)
	stableStore, err := boltdb.NewBoltStore(stableDB)
	if err != nil {
		return nil, fmt.Errorf("could not create stable store (path: %s): %w", stableDB, err)
	}

	// Create snapshot store. We never really expect we'll need snapshots
	// since our clusters are short lived, so this should be fine.
	snapshot := raft.NewDiscardSnapshotStore()

	fsm := newFsmExecutor(log, executor, cfg.Callbacks...)

	raftCfg := getRaftConfig(cfg, log, host.ID().String())

	// Tag the logger with the cluster ID (request ID).
	raftCfg.Logger = raftCfg.Logger.With("cluster", requestID)

	raftNode, err := raft.NewRaft(&raftCfg, fsm, logStore, stableStore, snapshot, transport)
	if err != nil {
		return nil, fmt.Errorf("could not create a raft node: %w", err)
	}

	rh := Handler{
		Raft:     raftNode,
		logStore: logStore,
		stable:   stableStore,

		log:     log.With().Str("module", "raft").Str("cluster", requestID).Logger(),
		cfg:     cfg,
		rootDir: rootDir,
		peers:   peers,
	}

	rh.log.Info().Strs("peers", blockless.PeerIDsToStr(peers)).Msg("created new raft handler")

	return &rh, nil
}

func (h *Handler) Shutdown() error {

	future := h.Raft.Shutdown()
	err := future.Error()
	if err != nil {
		return fmt.Errorf("could not shutdown raft cluster: %w", err)
	}

	// We'll log the actual error but return an "umbrella" one if we fail to close any of the two stores.
	var retErr error
	err = h.logStore.Close()
	if err != nil {
		log.Error().Err(err).Msg("could not close log store")
		retErr = fmt.Errorf("could not close raft database")
	}

	err = h.stable.Close()
	if err != nil {
		log.Error().Err(err).Msg("could not close stable store")
		retErr = fmt.Errorf("could not close raft database")
	}

	// Delete residual files. This may fail if we failed to close the databases above.
	err = os.RemoveAll(h.rootDir)
	if err != nil {
		log.Error().Err(err).Str("path", h.rootDir).Msg("could not delete consensus dir")
		retErr = fmt.Errorf("could not delete consensus directory")
	}

	return retErr
}

func (h *Handler) IsLeader() bool {
	return h.State() == raft.Leader
}
