package raft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/rs/zerolog"

	libp2praft "github.com/libp2p/go-libp2p-raft"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/host"
	"github.com/blessnetwork/b7s/models/blockless"
)

type Replica struct {
	*raft.Raft
	logStore *boltdb.BoltStore
	stable   *boltdb.BoltStore

	cfg Config
	log zerolog.Logger

	rootDir string
	peers   []peer.ID
}

// New creates a new raft replica, bootstraps the cluster and waits until a first leader is elected. We do this because
// only after the election the cluster is really operational and ready to process requests.
func New(log zerolog.Logger, host *host.Host, workspace string, requestID string, executor blockless.Executor, peers []peer.ID, options ...Option) (*Replica, error) {

	// Step 1: Create a new raft replica.
	replica, err := newReplica(log, host, workspace, requestID, executor, peers, options...)
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
			replica.log.Error().Type("type", obs.Data).Msg("invalid observation type received")
			return
		}

		// We don't need the observer anymore.
		replica.DeregisterObserver(observer)

		replica.log.Info().Str("request", requestID).Str("leader", string(leaderObs.LeaderID)).Msg("observed a leadership event - ready")
	}()

	replica.RegisterObserver(observer)

	// Step 3: Bootstrap the cluster.
	err = replica.bootstrapCluster()
	if err != nil {
		return nil, fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	wg.Wait()

	return replica, nil
}

func (r *Replica) Consensus() consensus.Type {
	return consensus.Raft
}

func newReplica(log zerolog.Logger, host *host.Host, workspace string, requestID string, executor blockless.Executor, peers []peer.ID, options ...Option) (*Replica, error) {

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

	rh := Replica{
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

func (r *Replica) Shutdown() error {

	r.log.Info().Msg("shuttting down cluster")

	future := r.Raft.Shutdown()
	err := future.Error()
	if err != nil {
		return fmt.Errorf("could not shutdown raft cluster: %w", err)
	}

	// We'll log the actual error but return an "umbrella" one if we fail to close any of the two stores.
	var multierr *multierror.Error

	err = r.logStore.Close()
	if err != nil {
		multierr = multierror.Append(multierr, fmt.Errorf("could not close log store: %w", err))
	}

	err = r.stable.Close()
	if err != nil {
		multierr = multierror.Append(multierr, fmt.Errorf("could not close stable store: %w", err))
	}

	// Delete residual files. This may fail if we failed to close the databases above.
	err = os.RemoveAll(r.rootDir)
	if err != nil {
		multierr = multierror.Append(multierr, fmt.Errorf("could not delete consensus dir: %w", err))
	}

	return multierr.ErrorOrNil()
}

func (r *Replica) isLeader() bool {
	return r.State() == raft.Leader
}
