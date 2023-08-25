package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/consensus"
	"github.com/blocklessnetworking/b7s/consensus/pbft"
	"github.com/blocklessnetworking/b7s/consensus/raft"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

// consensusExecutor defines the interface we have for managing clustered execution.
// Execute often does not mean a direct execution but instead just pipelining the request, where execution is done asynchronously.
type consensusExecutor interface {
	Consensus() consensus.Type
	Execute(from peer.ID, id string, request execute.Request) (codes.Code, execute.Result, error)
	Shutdown() error
}

func (n *Node) createRaftCluster(ctx context.Context, from peer.ID, fc request.FormCluster) error {

	// Add a callback function to cache the execution result
	cacheFn := func(req raft.FSMLogEntry, res execute.Result) {
		n.executeResponses.Set(req.RequestID, res)
	}

	// Add a callback function to send the execution result to origin.
	sendFn := func(req raft.FSMLogEntry, res execute.Result) {

		ctx, cancel := context.WithTimeout(context.Background(), consensusClusterSendTimeout)
		defer cancel()

		msg := response.Execute{
			Type:      blockless.MessageExecuteResponse,
			Code:      res.Code,
			RequestID: req.RequestID,
			Results: execute.ResultMap{
				n.host.ID(): res,
			},
		}

		err := n.send(ctx, req.Origin, msg)
		if err != nil {
			n.log.Error().Err(err).Str("peer", req.Origin.String()).Msg("could not send execution result to node")
		}
	}

	rh, err := raft.New(
		n.log,
		n.host,
		n.cfg.Workspace,
		fc.RequestID,
		n.executor,
		fc.Peers,
		raft.WithCallbacks(cacheFn, sendFn),
	)
	if err != nil {
		return fmt.Errorf("could not create raft node: %w", err)
	}

	n.clusterLock.Lock()
	n.clusters[fc.RequestID] = rh
	n.clusterLock.Unlock()

	res := response.FormCluster{
		Type:      blockless.MessageFormClusterResponse,
		RequestID: fc.RequestID,
		Code:      codes.OK,
		Consensus: fc.Consensus,
	}

	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send cluster confirmation message: %w", err)
	}

	return nil
}

func (n *Node) createPBFTCluster(ctx context.Context, from peer.ID, fc request.FormCluster) error {

	// TODO (pbft): Use an actual key, but we don't have signing yet.
	var dummyKey crypto.PrivKey
	ph, err := pbft.NewReplica(
		n.log,
		n.host,
		n.executor,
		fc.Peers,
		fc.RequestID,
		dummyKey,
	)
	if err != nil {
		return fmt.Errorf("could not create PBFT node: %w", err)
	}

	n.clusterLock.Lock()
	n.clusters[fc.RequestID] = ph
	n.clusterLock.Unlock()

	res := response.FormCluster{
		Type:      blockless.MessageFormClusterResponse,
		RequestID: fc.RequestID,
		Code:      codes.OK,
		Consensus: fc.Consensus,
	}

	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send cluster confirmation message: %w", err)
	}

	return nil
}

func (n *Node) leaveCluster(requestID string) error {

	// Shutdown can take a while so use short locking intervals.
	n.clusterLock.RLock()
	cluster, ok := n.clusters[requestID]
	n.clusterLock.RUnlock()

	if !ok {
		return errors.New("no cluster with that ID")
	}

	n.log.Info().Str("consensus", cluster.Consensus().String()).Str("request", requestID).Msg("leaving consensus cluster")

	ctx, cancel := context.WithTimeout(context.Background(), consensusClusterDisbandTimeout)
	defer cancel()

	// We know that the request is done executing when we have a result for it.
	_, ok = n.executeResponses.WaitFor(ctx, requestID)

	log := n.log.With().Str("request", requestID).Logger()
	log.Info().Bool("executed_work", ok).Msg("waiting for execution done, leaving cluster")

	err := cluster.Shutdown()
	if err != nil {
		// NOTE: Not much we can do at this point.
		return fmt.Errorf("could not leave cluster (request: %v): %w", requestID, err)
	}

	n.clusterLock.Lock()
	delete(n.clusters, requestID)
	n.clusterLock.Unlock()

	return nil
}

// helper function just for the sake of readibility.
func consensusRequired(c consensus.Type) bool {
	return c != 0
}
