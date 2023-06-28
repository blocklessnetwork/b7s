package node

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb/v2"

	libp2praft "github.com/libp2p/go-libp2p-raft"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/execute"
)

type raftHandler struct {
	*raft.Raft

	log    *boltdb.BoltStore
	stable *boltdb.BoltStore
}

func (n *Node) newRaftHandler(requestID string) (*raftHandler, error) {

	// Determine directory that should be used for consensus for this request.
	dirPath := n.consensusDir(requestID)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create consensus work directory: %w", err)
	}

	// Transport layer for raft communication.
	transport, err := libp2praft.NewLibp2pTransport(n.host, consensusTransportTimeout)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p transport: %w", err)
	}

	// Create log store.
	logDB := filepath.Join(dirPath, defaultLogStoreName)
	logStore, err := boltdb.NewBoltStore(logDB)
	if err != nil {
		return nil, fmt.Errorf("could not create log store (path: %s): %w", logDB, err)
	}

	// Create stable store.
	stableDB := filepath.Join(dirPath, defaultStableStoreName)
	stableStore, err := boltdb.NewBoltStore(stableDB)
	if err != nil {
		return nil, fmt.Errorf("could not create stable store (path: %s): %w", stableDB, err)
	}

	// Create snapshot store. We never really expect we'll need snapshots
	// since our clusters are short lived, so this should be fine.
	snapshot := raft.NewDiscardSnapshotStore()

	// Add a callback function to cache the execution result
	cacheFn := func(req fsmLogEntry, res execute.Result) {
		n.executeResponses.Set(req.RequestID, res)
	}

	fsm := newFsmExecutor(n.log, n.executor, cacheFn)

	raftCfg := n.getRaftConfig(n.host.ID().String())
	raftNode, err := raft.NewRaft(&raftCfg, fsm, logStore, stableStore, snapshot, transport)
	if err != nil {
		return nil, fmt.Errorf("could not create a raft node: %w", err)
	}

	rh := raftHandler{
		Raft:   raftNode,
		log:    logStore,
		stable: stableStore,
	}

	return &rh, nil
}

func (n *Node) getRaftConfig(nodeID string) raft.Config {
	// TODO: (raft): use zerolog here, not a random hclog instance, even if it is JSON.
	logOpts := hclog.LoggerOptions{
		JSONFormat: true,
		Level:      hclog.Debug,
		Output:     os.Stderr,
		Name:       "raft",
	}
	raftLogger := hclog.New(&logOpts)

	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(nodeID)
	cfg.Logger = raftLogger
	cfg.HeartbeatTimeout = n.cfg.ConsensusHeartbeatTimeout
	cfg.ElectionTimeout = n.cfg.ConsensusElectionTimeout
	cfg.LeaderLeaseTimeout = n.cfg.ConsensusLeaderLease

	return *cfg
}

func bootstrapCluster(raftHandler *raftHandler, peerIDs []peer.ID) error {

	if len(peerIDs) == 0 {
		return errors.New("empty peer list")
	}

	servers := make([]raft.Server, 0, len(peerIDs))
	for _, id := range peerIDs {

		s := raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(id.String()),
			Address:  raft.ServerAddress(id),
		}

		servers = append(servers, s)
	}

	cfg := raft.Configuration{
		Servers: servers,
	}

	// Bootstrapping will only succeed for the first node to start it.
	// Other attempts will fail with an error that can be ignored.
	ret := raftHandler.BootstrapCluster(cfg)
	err := ret.Error()
	if err != nil && !errors.Is(err, raft.ErrCantBootstrap) {
		return fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	return nil
}

func (n *Node) shutdownCluster(requestID string) error {

	n.log.Info().Str("request_id", requestID).Msg("shutting down cluster")

	n.clusterLock.RLock()
	raftHandler, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	if !ok {
		return nil
	}

	future := raftHandler.Shutdown()
	err := future.Error()
	if err != nil {
		return fmt.Errorf("could not shutdown raft cluster: %w", err)
	}

	// We'll log the actual error but return an "umbrella" one if we fail to close any of the two stores.
	var retErr error
	err = raftHandler.log.Close()
	if err != nil {
		n.log.Error().Err(err).Str("request_id", requestID).Msg("could not close log store")
		retErr = fmt.Errorf("could not close raft database")
	}

	err = raftHandler.stable.Close()
	if err != nil {
		n.log.Error().Err(err).Str("request_id", requestID).Msg("could not close stable store")
		retErr = fmt.Errorf("could not close raft database")
	}

	// Delete residual files. This may fail if we failed to close the databases above.
	dir := n.consensusDir(requestID)
	err = os.RemoveAll(dir)
	if err != nil {
		n.log.Error().Err(err).Str("request_id", requestID).Str("path", dir).Msg("could not delete consensus dir")
		retErr = fmt.Errorf("could not delete consensus directory")
	}

	n.clusterLock.Lock()
	delete(n.clusters, requestID)
	n.clusterLock.Unlock()

	return retErr
}

func (n *Node) consensusDir(requestID string) string {
	return filepath.Join(n.cfg.Workspace, defaultConsensusDirName, requestID)
}
